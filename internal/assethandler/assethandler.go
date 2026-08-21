package assethandler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/bragemusic/core/pkg/config"
	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/types"
)

const disableTransitions = `/* Disable Transitions */
.tooltip[data-entering],
.tooltip[data-exiting],
.select__popover[data-entering],
.select__popover[data-exiting],
.dropdown__popover[data-entering],
.dropdown__popover[data-exiting],
.modal__container[data-entering],
.modal__container[data-exiting],
.popover[data-entering],
.popover[data-exiting] {
  animation: none !important;
  transition: none !important;
}
`

type AssetHandler struct {
	ClientType         types.DeviceType
	ServerBaseURL      string
	Themes             map[string]types.Theme
	DisableTransitions bool
	ImgFolderPath      string
	sc                 *serverclient.ServerClient
}

func (a *AssetHandler) generateCustomCSS() string {
	tcss := []string{}

	if a.DisableTransitions {
		tcss = append(tcss, disableTransitions)
	}

	for themeName, theme := range a.Themes {
		tcss = append(tcss, theme.CSS(themeName))
	}

	return strings.Join(tcss, "\n")
}

func (a *AssetHandler) getLocalImage(w http.ResponseWriter, r *http.Request) {
	requestedFilename := strings.TrimPrefix(r.URL.Path, "/api/img/")
	fileData, err := os.ReadFile(path.Join(a.ImgFolderPath, requestedFilename))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Could not load file %s", requestedFilename)
	}

	w.Write(fileData)
}

func (a *AssetHandler) getServerImage(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", a.ServerBaseURL, r.URL.Path), nil)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.sc.AuthToken()))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "upstream error", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// copy headers (important for caching & content-type)
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("Cache-Control", resp.Header.Get("Cache-Control"))

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (a *AssetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.String() {
	case "/config/custom.css":
		w.Header().Set("Content-Type", "text/css; charset=utf-8")

		_, err := w.Write([]byte(a.generateCustomCSS()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if strings.HasPrefix(r.URL.Path, "/api/img") {
		if strings.HasPrefix(r.URL.Path, "/api/img/users") {
			a.getServerImage(w, r)
			return
		}

		switch a.ClientType {
		case types.DeviceTypeStreaming:
			a.getServerImage(w, r)
			return
		case types.DeviceTypeSync:
			a.getLocalImage(w, r)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func New(cfg config.ClientConfig, sc *serverclient.ServerClient) AssetHandler {
	return AssetHandler{
		ServerBaseURL:      cfg.Server.BaseUrl,
		ClientType:         cfg.General.ClientType,
		Themes:             cfg.Themes,
		DisableTransitions: cfg.General.DisableTransitions,
		sc:                 sc,
		ImgFolderPath:      cfg.Paths.ImageDir,
	}
}
