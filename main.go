package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/alexzorin/lpass-ui/lpass"
	"github.com/atotto/clipboard"
	"github.com/pkg/errors"
	"gopkg.in/qml.v1"
)

var (
	queryMu     sync.Mutex
	prevQuery   string // last query to be processed
	latestQuery string // the latest query
	root        *Control
)

type Control struct {
	Root qml.Object
}

func main() {
	if err := qml.Run(run); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting ui: %v\n", err)
		os.Exit(1)
	}
}

// UI entry point
func run() error {
	go query()

	engine := qml.NewEngine()

	controls, err := engine.Load("lpass-ui.qml", strings.NewReader(mainQml))
	if err != nil {
		return errors.Wrap(err, "Failed to load qml file")
	}

	window := controls.CreateWindow(nil)
	root = &Control{Root: window.Root()}

	context := engine.Context()
	context.SetVar("ctrl", root)

	window.Set("flags", 0x00000002|0x00000001) // Qt::Dialog
	window.Set("title", "LastPass UI")
	window.Show()

	window.Wait()

	return nil
}

func query() {
	for _ = range time.Tick(250 * time.Millisecond) {
		queryMu.Lock()
		q := strings.TrimSpace(latestQuery)
		queryMu.Unlock()
		if prevQuery == q || len(q) < 3 {
			continue
		}

		res, err := lpass.QuerySites(q)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to retreive sites: %v\n", err)
		}

		m := root.Root.ObjectByName("resultsModel")
		m.Call("clear")
		for _, v := range res {
			m.Call("addResult", v.Name, v.Username, v.Password)
		}

		prevQuery = q
	}
}

// When we are typing chars into the query
func (c *Control) OnQueryInput(text qml.Object) {
	queryMu.Lock()
	latestQuery = text.String("text")
	queryMu.Unlock()
}

// Enter is pressed on an item - password is copied
func (c *Control) OnAccepted(m qml.Object) {
	clipboard.WriteAll(m.String("password"))
	os.Exit(0)
}

func (c *Control) Bail() {
	os.Exit(0)
}

const mainQml = `import QtQuick 2.3
import QtQuick.Layouts 1.2
import QtQuick.Controls 1.4


Rectangle {
  id: root
  width: 600
  height: 200
  color: 'white'

  ColumnLayout {
    spacing: 2

    // Where we take our input
    TextInput {
      id: queryInput
      Layout.alignment: Qt.AlignCenter
      Layout.preferredWidth: 600
      Layout.preferredHeight: 50
      width: 600
      height: parent.height / 2
      color: "#333"
      font.pixelSize: 24
      wrapMode: TextInput.Wrap
      focus: true
      Keys.onReleased: ctrl.onQueryInput(queryInput)
      Keys.onTabPressed: {
        var count = resultsTable.rowCount;
	var cur = resultsTable.currentRow;
	if(cur < 0) {
	  cur = 0;
	} else if (cur == count-1) {
          cur = 0;
	} else {
	  cur++;
	}
	resultsTable.currentRow = cur;
	resultsTable.selection.clear();
	resultsTable.selection.select(cur);
      }
      Keys.onEscapePressed: ctrl.bail()
      onAccepted: { 
        if(resultsTable.rowCount > 0 && resultsTable.currentRow != -1) { 
       	  ctrl.onAccepted(resultsModel.get(resultsTable.currentRow))
	}
      }
    }

    // Where we list the matches
    TableView {
      id: resultsTable
      Layout.alignment: Qt.AlignCenter
      Layout.preferredWidth: 600
      Layout.preferredHeight: 150
      selectionMode: SelectionMode.SingleSelection

      model: resultsModel

      ListModel {
        id: resultsModel
	objectName: "resultsModel"
	function addResult(name, username, password) {
	  resultsModel.append({"name": name, "username": username, "password": password});
	}
      }

      TableViewColumn {
        role: "name"
	title: "Name"
	width: 300 
      }
      TableViewColumn {
        role: "username"
	title: "Username"
	width: 300
      }
    }

  }

}

`
