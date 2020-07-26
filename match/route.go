// Go HTTP router based on a simple custom match() function

package match

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func Route(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	var slug string
	var id int

	p := r.URL.Path
	switch {
	case match(p, "/"):
		h = get(home)
	case match(p, "/contact"):
		h = get(contact)
	case match(p, "/api/widgets") && r.Method == "GET":
		h = get(apiGetWidgets)
	case match(p, "/api/widgets"):
		h = post(apiCreateWidget)
	case match(p, "/api/widgets/+", &slug):
		h = post(apiWidget{slug}.update)
	case match(p, "/api/widgets/+/parts", &slug):
		h = post(apiWidget{slug}.createPart)
	case match(p, "/api/widgets/+/parts/+/update", &slug, &id):
		h = post(apiWidgetPart{slug, id}.update)
	case match(p, "/api/widgets/+/parts/+/delete", &slug, &id):
		h = post(apiWidgetPart{slug, id}.delete)
	case match(p, "/+", &slug):
		h = get(widget{slug}.widget)
	case match(p, "/+/admin", &slug):
		h = get(widget{slug}.admin)
	case match(p, "/+/image", &slug):
		h = post(widget{slug}.image)
	default:
		http.NotFound(w, r)
		return
	}
	h.ServeHTTP(w, r)
}

func match(path, pattern string, vars ...interface{}) bool {
	for ; pattern != "" && path != ""; pattern = pattern[1:] {
		switch pattern[0] {
		case '+':
			// '+' matches till next slash in path
			slash := strings.IndexByte(path, '/')
			if slash < 0 {
				slash = len(path)
			}
			segment := path[:slash]
			path = path[slash:]
			switch p := vars[0].(type) {
			case *string:
				*p = segment
			case *int:
				n, err := strconv.Atoi(segment)
				if err != nil {
					return false
				}
				*p = n
			default:
				panic("vars must be *string or *int")
			}
			vars = vars[1:]
		case path[0]:
			// non-'+' pattern byte must match path byte
			path = path[1:]
		default:
			return false
		}
	}
	return path == "" && pattern == ""
}

func allowMethod(h http.HandlerFunc, methods ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, m := range methods {
			if m == r.Method {
				h(w, r)
				return
			}
		}
		w.Header().Set("Allow", strings.Join(methods, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
	}
}

func get(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, http.MethodGet, http.MethodHead)
}

func post(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, http.MethodPost)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "home\n")
}

func contact(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "contact\n")
}

func apiGetWidgets(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "apiGetWidgets\n")
}

func apiCreateWidget(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "apiCreateWidget\n")
}

type apiWidget struct {
	slug string
}

func (h apiWidget) update(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiUpdateWidget %s\n", h.slug)
}

func (h apiWidget) createPart(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiCreateWidgetPart %s\n", h.slug)
}

type apiWidgetPart struct {
	slug string
	id   int
}

func (h apiWidgetPart) update(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiUpdateWidgetPart %s %d\n", h.slug, h.id)
}

func (h apiWidgetPart) delete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiDeleteWidgetPart %s %d\n", h.slug, h.id)
}

type widget struct {
	slug string
}

func (h widget) widget(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "widget %s\n", h.slug)
}

func (h widget) admin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "widgetAdmin %s\n", h.slug)
}

func (h widget) image(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "widgetImage %s\n", h.slug)
}