package output

import (
	"net/http"

	"bytes"
	"fmt"
	"io"
	"os/exec"

	"github.com/barasher/dep-carto/internal/model"
)

type JpgFormatter struct{}

func NewJpgFormatter() JpgFormatter {
	return JpgFormatter{}
}

func (JpgFormatter) Format(s []model.Server, w http.ResponseWriter) {
	df := NewDotFormatter()
	b := bytes.Buffer{}
	if err := df.FormatToWriter(s, &b); err != nil {
		http.Error(w, "error while converting servers to dot: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	if err := dotToJpg(&b, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func dotToJpg(r io.Reader, w io.Writer) error {
	cmd := exec.Command("dot", "-Tjpg")
	b := bytes.Buffer{}
	cmd.Stdin = r
	cmd.Stdout = w
	cmd.Stderr = &b
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error while converting: %v", b.String())
	}
	return nil
}
