package templates

import (
	"context"
	"strings"

	"github.com/a-h/templ"
)

func RenderToString(ctx context.Context, component templ.Component) string {
	var w strings.Builder
	component.Render(ctx, &w)
	return w.String()
}
