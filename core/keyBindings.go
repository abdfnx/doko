package core

import (
	"io"
	"os"
	"fmt"
	"errors"
	"context"
	"os/signal"
	"path/filepath"

	"github.com/abdfnx/doko/log"
	"github.com/abdfnx/doko/docker"
	"github.com/abdfnx/doko/shared"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
)

var inputWidth = 70

func (ui *UI) setGlobalKeybinding(event *tcell.EventKey) {
	switch event.Rune() {
		case 'h':
			ui.prevPanel()
		case 'l':
			ui.nextPanel()
		case 'q':
			ui.Stop()
		case '/':
			ui.filter()
	}

	switch event.Key() {
		case tcell.KeyTab:
			ui.nextPanel()
		case tcell.KeyBacktab:
			ui.prevPanel()
		case tcell.KeyRight:
			ui.nextPanel()
		case tcell.KeyLeft:
			ui.prevPanel()
	}
}

func (ui *UI) filter() {
	currentPanel := ui.state.panels.panel[ui.state.panels.currentPanel]
	if currentPanel.name() == "tasks" {
		return
	}

	currentPanel.setFilterWord("")
	currentPanel.updateEntries(ui)

	viewName := "filter"
	searchInput := tview.NewInputField().SetLabel("Word")
	searchInput.SetLabelWidth(6)
	searchInput.SetTitle("filter")
	searchInput.SetTitleAlign(tview.AlignLeft)
	searchInput.SetBorder(true)

	closeSearchInput := func() {
		ui.closeAndSwitchPanel(viewName, ui.state.panels.panel[ui.state.panels.currentPanel].name())
	}

	searchInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			closeSearchInput()
		}
	})

	searchInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			closeSearchInput()
		}
		return event
	})

	searchInput.SetChangedFunc(func(text string) {
		currentPanel.setFilterWord(text)
		currentPanel.updateEntries(ui)
	})

	ui.pages.AddAndSwitchToPage(viewName, ui.modal(searchInput, 80, 3), true).ShowPage("main")
}

func (ui *UI) nextPanel() {
	idx := (ui.state.panels.currentPanel + 1) % len(ui.state.panels.panel)
	ui.switchPanel(ui.state.panels.panel[idx].name())
}

func (ui *UI) prevPanel() {
	ui.state.panels.currentPanel--

	if ui.state.panels.currentPanel < 0 {
		ui.state.panels.currentPanel = len(ui.state.panels.panel) - 1
	}

	idx := (ui.state.panels.currentPanel) % len(ui.state.panels.panel)
	ui.switchPanel(ui.state.panels.panel[idx].name())
}

func (ui *UI) createContainerForm() {
	selectedImage := ui.selectedImage()

	if selectedImage == nil {
		logger.Logger.Error("please input image")
		return
	}

	image := fmt.Sprintf("%s:%s", selectedImage.Repo, selectedImage.Tag)

	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle("Create container")
	form.SetTitleAlign(tview.AlignLeft)

	form.AddInputField("Name", "", inputWidth, nil, nil).
		AddInputField("HostIP", "", inputWidth, nil, nil).
		AddInputField("HostPort", "", inputWidth, nil, nil).
		AddInputField("Port", "", inputWidth, nil, nil).
		AddDropDown("VolumeType", []string{"bind", "volume"}, 0, func(option string, optionIndex int) {}).
		AddInputField("HostVolume", "", inputWidth, nil, nil).
		AddInputField("Volume", "", inputWidth, nil, nil).
		AddInputField("Image", image, inputWidth, nil, nil).
		AddInputField("User", "", inputWidth, nil, nil).
		AddCheckbox("Attach", false, nil).
		AddInputField("Env", "", inputWidth, nil, nil).
		AddInputField("Cmd", "", inputWidth, nil, nil).
		AddButton("Create", func() {
			ui.createContainer(form, image)
		}).
		AddButton("Cancel", func() {
			ui.closeAndSwitchPanel("form", "images")
		})

	ui.pages.AddAndSwitchToPage("form", ui.modal(form, 80, 29), true).ShowPage("main")
}

func (ui *UI) createContainer(form *tview.Form, image string) {
	ui.startTask("create container "+image, func(ctx context.Context) error {
		inputLabels := []string{
			"Name",
			"HostIP",
			"Port",
			"HostVolume",
			"Volume",
			"Image",
			"User",
		}

		var data = make(map[string]string)

		for _, label := range inputLabels {
			data[label] = form.GetFormItemByLabel(label).(*tview.InputField).GetText()
		}

		_, volumeType := form.GetFormItemByLabel("VolumeType").(*tview.DropDown).
			GetCurrentOption()
		data["VolumeType"] = volumeType

		isAttach := form.GetFormItemByLabel("Attach").(*tview.Checkbox).IsChecked()

		options, err := docker.Client.NewContainerOptions(data, isAttach)
		if err != nil {
			logger.Logger.Errorf("cannot create container %s", err)
			return err
		}

		err = docker.Client.CreateContainer(options)

		if err != nil {
			logger.Logger.Errorf("cannot create container %s", err)
			return err
		}

		ui.closeAndSwitchPanel("form", "images")
		go ui.app.QueueUpdateDraw(func() {
			ui.containerPanel().setEntries(ui)
		})

		return nil
	})
}

func (ui *UI) pullImageForm() {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitleAlign(tview.AlignLeft)
	form.SetTitle("Pull image")
	form.AddInputField("Image", "", inputWidth, nil, nil).
		AddButton("Pull", func() {
			image := form.GetFormItemByLabel("Image").(*tview.InputField).GetText()
			ui.pullImage(image, "form", "images")
		}).
		AddButton("Cancel", func() {
			ui.closeAndSwitchPanel("form", "images")
		})

	ui.pages.AddAndSwitchToPage("form", ui.modal(form, 80, 7), true).ShowPage("main")
}

func (ui *UI) pullImage(image, closePanel, switchPanel string) {
	ui.startTask("Pull image "+image, func(ctx context.Context) error {
		ui.closeAndSwitchPanel(closePanel, switchPanel)
		err := docker.Client.PullImage(image)
		if err != nil {
			logger.Logger.Errorf("cannot pull an image %s", err)
			return err
		}

		ui.imagePanel().updateEntries(ui)

		return nil
	})
}

func (ui *UI) displayInspect(data, page string) {
	text := tview.NewTextView()
	text.SetTitle("Detail").SetTitleAlign(tview.AlignLeft)
	text.SetBorder(true)
	text.SetText(data)

	text.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Rune() == 'q' {
			ui.closeAndSwitchPanel("detail", page)
		}

		return event
	})

	ui.pages.AddAndSwitchToPage("detail", text, true)
}

func (ui *UI) inspectImage() {
	image := ui.selectedImage()

	inspect, err := docker.Client.InspectImage(image.ID)

	if err != nil {
		logger.Logger.Errorf("cannot inspect image %s", err)
		return
	}

	ui.displayInspect(shared.StructToJSON(inspect), "images")
}

func (ui *UI) renameContainerForm() {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitleAlign(tview.AlignLeft)
	form.SetTitle("Rename container")
	form.AddInputField("NewName", "", inputWidth, nil, nil).
		AddButton("Rename", func() {
			image := form.GetFormItemByLabel("NewName").(*tview.InputField).GetText()
			ui.renameContainer(image, "form", "containers")
		}).
		AddButton("Cancel", func() {
			ui.closeAndSwitchPanel("form", "containers")
		})

	ui.pages.AddAndSwitchToPage("form", ui.modal(form, 80, 7), true).ShowPage("main")
}

func (ui *UI) renameContainer(newName, closePanel, switchPanel string) {
	ui.startTask("Renaming container "+newName, func(ctx context.Context) error {
		ui.closeAndSwitchPanel(closePanel, switchPanel)
		oldContainer := ui.selectedContainer()

		if oldContainer == nil {
			err := errors.New("specified container is nil")
			logger.Logger.Errorf("cannot rename container %s", err)
			return err
		}

		err := docker.Client.RenameContainer(oldContainer.ID, newName)
		if err != nil {
			logger.Logger.Errorf("cannot create container %s", err)
			return err
		}

		ui.containerPanel().updateEntries(ui)

		return nil
	})
}

func (ui *UI) inspectContainer() {
	container := ui.selectedContainer()

	inspect, err := docker.Client.InspectContainer(container.ID)
	if err != nil {
		logger.Logger.Errorf("cannot inspect container %s", err)
		return
	}

	ui.displayInspect(shared.StructToJSON(inspect), "containers")
}

func (ui *UI) inspectVolume() {
	volume := ui.selectedVolume()

	inspect, err := docker.Client.InspectVolume(volume.Name)

	if err != nil {
		logger.Logger.Errorf("cannot inspect volume %s", err)
		return
	}

	ui.displayInspect(shared.StructToJSON(inspect), "volumes")
}

func (ui *UI) inspectNetwork() {
	network := ui.selectedNetwork()

	inspect, err := docker.Client.InspectNetwork(network.ID)
	if err != nil {
		logger.Logger.Errorf("cannot inspect network %s", err)
		return
	}

	ui.displayInspect(shared.StructToJSON(inspect), "networks")
}

func (ui *UI) removeImage() {
	image := ui.selectedImage()

	ui.confirm("Do you want to remove the image?", "Done", "images", func() {
		ui.startTask(fmt.Sprintf("remove image %s:%s", image.Repo, image.Tag), func(ctx context.Context) error {
			if err := docker.Client.RemoveImage(image.ID); err != nil {
				logger.Logger.Errorf("cannot remove the image %s", err)
				return err
			}
			ui.imagePanel().updateEntries(ui)
			return nil
		})
	})
}

func (ui *UI) removeContainer() {
	container := ui.selectedContainer()

	ui.confirm("Do you want to remove the container?", "Done", "containers", func() {
		ui.startTask(fmt.Sprintf("remove container %s", container.Name), func(ctx context.Context) error {
			if err := docker.Client.RemoveContainer(container.ID); err != nil {
				logger.Logger.Errorf("cannot remove the container %s", err)
				return err
			}
			ui.containerPanel().updateEntries(ui)
			return nil
		})
	})
}

func (ui *UI) removeVolume() {
	volume := ui.selectedVolume()

	ui.confirm("Do you want to remove the volume?", "Done", "volumes", func() {
		ui.startTask(fmt.Sprintf("remove volume %s", volume.Name), func(ctx context.Context) error {
			if err := docker.Client.RemoveVolume(volume.Name); err != nil {
				logger.Logger.Errorf("cannot remove the volume %s", err)
				return err
			}
			ui.volumePanel().updateEntries(ui)
			return nil
		})
	})
}

func (ui *UI) removeNetwork() {
	network := ui.selectedNetwork()

	ui.confirm("Do you want to remove the network?", "Done", "networks", func() {
		ui.startTask(fmt.Sprintf("remove network %s", network.Name), func(ctx context.Context) error {
			if err := docker.Client.RemoveNetwork(network.ID); err != nil {
				logger.Logger.Errorf("cannot remove the network %s", err)
				return err
			}

			ui.networkPanel().updateEntries(ui)
			return nil
		})
	})
}

func (ui *UI) startContainer() {
	container := ui.selectedContainer()

	ui.startTask(fmt.Sprintf("start container %s", container.Name), func(ctx context.Context) error {
		if err := docker.Client.StartContainer(container.ID); err != nil {
			logger.Logger.Errorf("cannot start container %s", err)
			return err
		}

		ui.containerPanel().updateEntries(ui)
		return nil
	})
}

func (ui *UI) stopContainer() {
	container := ui.selectedContainer()

	ui.startTask(fmt.Sprintf("stop container %s", container.Name), func(ctx context.Context) error {

		if err := docker.Client.StopContainer(container.ID); err != nil {
			logger.Logger.Errorf("cannot stop container %s", err)
			return err
		}

		ui.containerPanel().updateEntries(ui)
		return nil
	})
}

func (ui *UI) exportContainerForm() {
	inputWidth := 70

	container := ui.selectedContainer()
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitleAlign(tview.AlignLeft)
	form.SetTitle("Export container")
	form.AddInputField("Path", "", inputWidth, nil, nil).
		AddInputField("Container", container.Name, inputWidth, nil, nil).
		AddButton("Create", func() {
			path := form.GetFormItemByLabel("Path").(*tview.InputField).GetText()
			container := form.GetFormItemByLabel("Container").(*tview.InputField).GetText()

			ui.exportContainer(path, container)
		}).
		AddButton("Cancel", func() {
			ui.closeAndSwitchPanel("form", "containers")
		})

	ui.pages.AddAndSwitchToPage("form", ui.modal(form, 80, 9), true).ShowPage("main")
}

func (ui *UI) exportContainer(path, container string) {
	ui.startTask("export container "+container, func(ctx context.Context) error {
		ui.closeAndSwitchPanel("form", "containers")
		err := docker.Client.ExportContainer(container, path)
		if err != nil {
			logger.Logger.Errorf("cannot export container %s", err)
			return err
		}

		return nil
	})
}

func (ui *UI) loadImageForm() {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitleAlign(tview.AlignLeft)
	form.SetTitle("Load image")
	form.AddInputField("Path", "", inputWidth, nil, nil).
		AddButton("Load", func() {
			path := form.GetFormItemByLabel("Path").(*tview.InputField).GetText()
			ui.loadImage(path)
		}).
		AddButton("Cancel", func() {
			ui.closeAndSwitchPanel("form", "images")
		})

	ui.pages.AddAndSwitchToPage("form", ui.modal(form, 80, 7), true).ShowPage("main")
}

func (ui *UI) loadImage(path string) {
	ui.startTask("load image "+filepath.Base(path), func(ctx context.Context) error {
		ui.closeAndSwitchPanel("form", "images")
		if err := docker.Client.LoadImage(path); err != nil {
			logger.Logger.Errorf("cannot load image %s", err)
			return err
		}

		ui.imagePanel().updateEntries(ui)
		return nil
	})
}

func (ui *UI) importImageForm() {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitleAlign(tview.AlignLeft)
	form.SetTitle("Import image")
	form.AddInputField("Repository", "", inputWidth, nil, nil).
		AddInputField("Tag", "", inputWidth, nil, nil).
		AddInputField("Path", "", inputWidth, nil, nil).
		AddButton("Load", func() {
			repository := form.GetFormItemByLabel("Repository").(*tview.InputField).GetText()
			tag := form.GetFormItemByLabel("Tag").(*tview.InputField).GetText()
			path := form.GetFormItemByLabel("Path").(*tview.InputField).GetText()
			ui.importImage(path, repository, tag)
		}).
		AddButton("Cancel", func() {
			ui.closeAndSwitchPanel("form", "images")
		})

	ui.pages.AddAndSwitchToPage("form", ui.modal(form, 80, 11), true).ShowPage("main")
}

func (ui *UI) importImage(file, repo, tag string) {
	ui.startTask("import image "+file, func(ctx context.Context) error {
		ui.closeAndSwitchPanel("form", "images")

		if err := docker.Client.ImportImage(repo, tag, file); err != nil {
			logger.Logger.Errorf("cannot load image %s", err)
			return err
		}

		ui.imagePanel().updateEntries(ui)
		return nil
	})
}

func (ui *UI) saveImageForm() {
	image := ui.selectedImage()
	imageName := fmt.Sprintf("%s:%s", image.Repo, image.Tag)

	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitleAlign(tview.AlignLeft)
	form.SetTitle("Save image")
	form.AddInputField("Path", "", inputWidth, nil, nil).
		AddInputField("Image", imageName, inputWidth, nil, nil).
		AddButton("Save", func() {
			image := form.GetFormItemByLabel("Image").(*tview.InputField).GetText()
			path := form.GetFormItemByLabel("Path").(*tview.InputField).GetText()
			ui.saveImage(image, path)
		}).
		AddButton("Cancel", func() {
			ui.closeAndSwitchPanel("form", "images")
		})

	ui.pages.AddAndSwitchToPage("form", ui.modal(form, 80, 9), true).ShowPage("main")

}

func (ui *UI) saveImage(image, path string) {
	ui.startTask("save image "+image, func(ctx context.Context) error {
		ui.closeAndSwitchPanel("form", "images")

		if err := docker.Client.SaveImage([]string{image}, path); err != nil {
			logger.Logger.Errorf("cannot save image %s", err)
			return err
		}

		return nil
	})

}

func (ui *UI) commitContainerForm() {
	container := ui.selectedContainer()

	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitleAlign(tview.AlignLeft)
	form.SetTitle("Commit container")
	form.AddInputField("Repository", "", inputWidth, nil, nil).
		AddInputField("Tag", "", inputWidth, nil, nil).
		AddInputField("Container", container.Name, inputWidth, nil, nil).
		AddButton("Commit", func() {
			repo := form.GetFormItemByLabel("Repository").(*tview.InputField).GetText()
			tag := form.GetFormItemByLabel("Tag").(*tview.InputField).GetText()
			con := form.GetFormItemByLabel("Container").(*tview.InputField).GetText()
			ui.commitContainer(repo, tag, con)
		}).
		AddButton("Cancel", func() {
			ui.closeAndSwitchPanel("form", "containers")
		})

	ui.pages.AddAndSwitchToPage("form", ui.modal(form, 80, 11), true).ShowPage("main")
}

func (ui *UI) commitContainer(repo, tag, container string) {
	ui.startTask("commit container "+container, func(ctx context.Context) error {
		ui.closeAndSwitchPanel("form", "containers")

		if err := docker.Client.CommitContainer(container, types.ContainerCommitOptions{Reference: repo + ":" + tag}); err != nil {
			logger.Logger.Errorf("cannot commit container %s", err)
			return err
		}

		ui.imagePanel().updateEntries(ui)
		return nil
	})
}

func (ui *UI) attachContainerForm() {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitleAlign(tview.AlignLeft)
	form.SetTitle("Exec container")
	form.AddInputField("Cmd", "", inputWidth, nil, nil).
		AddButton("Exec", func() {
			cmd := form.GetFormItemByLabel("Cmd").(*tview.InputField).GetText()
			ui.attachContainer(ui.selectedContainer().ID, cmd)
		}).
		AddButton("Cancel", func() {
			ui.closeAndSwitchPanel("form", "containers")
		})

	ui.pages.AddAndSwitchToPage("form", ui.modal(form, 80, 7), true).ShowPage("main")
}

func (ui *UI) attachContainer(container, cmd string) {
	ui.closeAndSwitchPanel("form", "containers")

	if !ui.app.Suspend(func() {
		ui.stopMonitoring()
		if err := docker.Client.AttachExecContainer(container, cmd); err != nil {
			logger.Logger.Errorf("cannot attach container %s", err)
		}

		ui.startMonitoring()
	}) {
		logger.Logger.Error("cannot suspend tview")
	}
}

func (ui *UI) createVolumeForm() {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitleAlign(tview.AlignLeft)
	form.SetTitle("Create volume")
	form.AddInputField("Name", "", inputWidth, nil, nil).
		AddInputField("Labels", "", inputWidth, nil, nil).
		AddInputField("Driver", "", inputWidth, nil, nil).
		AddInputField("Options", "", inputWidth, nil, nil).
		AddButton("Create", func() {
			ui.createVolume(form)
		}).
		AddButton("Cancel", func() {
			ui.closeAndSwitchPanel("form", "volumes")
		})

	ui.pages.AddAndSwitchToPage("form", ui.modal(form, 80, 13), true).ShowPage("main")
}

func (ui *UI) createVolume(form *tview.Form) {
	var data = make(map[string]string)
	inputLabels := []string{
		"Name",
		"Labels",
		"Driver",
		"Options",
	}

	for _, label := range inputLabels {
		data[label] = form.GetFormItemByLabel(label).(*tview.InputField).GetText()
	}

	ui.startTask("create volume "+data["Name"], func(ctx context.Context) error {
		options := docker.Client.NewCreateVolumeOptions(data)

		if err := docker.Client.CreateVolume(options); err != nil {
			logger.Logger.Errorf("cannot create volume %s", err)
			return err
		}

		ui.closeAndSwitchPanel("form", "volumes")

		go ui.app.QueueUpdateDraw(func() {
			ui.volumePanel().setEntries(ui)
		})

		return nil
	})
}

func (ui *UI) tailContainerLog() {
	container := ui.selectedContainer()

	if container == nil {
		logger.Logger.Errorf("cannot start tail container: selected container is null")
		return
	}

	if !ui.app.Suspend(func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		errCh := make(chan error)

		var reader io.ReadCloser
		var err error

		go func() {
			reader, err = docker.Client.ContainerLogStream(container.ID)
			if err != nil {
				logger.Logger.Error(err)
				errCh <- err
			}

			defer reader.Close()

			_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, reader)
			if err != nil {
				logger.Logger.Error(err)
				errCh <- err
			}

			return
		}()

		select {
			case err := <-errCh:
				logger.Logger.Error(err)
				reader.Close()
				return
			case <-sigint:
				reader.Close()
				return
		}
	}) {
		logger.Logger.Error("cannot suspend tview")
	}
}

func (ui *UI) killContainer() {
	container := ui.selectedContainer()

	if container == nil {
		logger.Logger.Errorf("cannot kill container: selected container is null")
		return
	}

	ui.confirm("Do you want to kill the container?", "Done", "containers", func() {
		ui.startTask(fmt.Sprintf("kill container %s", container.Name), func(ctx context.Context) error {
			if err := docker.Client.KillContainer(container.ID); err != nil {
				logger.Logger.Errorf("cannot kill the container %s", err)
				return err
			}

			ui.containerPanel().updateEntries(ui)
			return nil
		})
	})
}
