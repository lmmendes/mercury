package api

import (
	"mercury/internal/core"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	core   *core.Core
	router *mux.Router
}

func NewServer(core *core.Core) *Server {
	s := &Server{
		core:   core,
		router: mux.NewRouter(),
	}

	s.routes()
	return s
}

func (s *Server) withRecovery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.core.Logger.Printf("panic recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

func (s *Server) routes() {
	s.router.HandleFunc("/accounts", s.withRecovery(s.createAccount)).Methods("POST")
	s.router.HandleFunc("/accounts", s.withRecovery(s.getAccounts)).Methods("GET")
	s.router.HandleFunc("/accounts/{id}", s.withRecovery(s.getAccount)).Methods("GET")
	s.router.HandleFunc("/accounts/{id}", s.withRecovery(s.updateAccount)).Methods("PUT")
	s.router.HandleFunc("/accounts/{id}", s.withRecovery(s.deleteAccount)).Methods("DELETE")

	s.router.HandleFunc("/accounts/{accountId}/inboxes", s.withRecovery(s.createInbox)).Methods("POST")
	s.router.HandleFunc("/accounts/{accountId}/inboxes", s.withRecovery(s.getInboxes)).Methods("GET")
	s.router.HandleFunc("/accounts/{accountId}/inboxes/{inboxId}", s.withRecovery(s.getInbox)).Methods("GET")
	s.router.HandleFunc("/accounts/{accountId}/inboxes/{inboxId}", s.withRecovery(s.updateInbox)).Methods("PUT")
	s.router.HandleFunc("/accounts/{accountId}/inboxes/{inboxId}", s.withRecovery(s.deleteInbox)).Methods("DELETE")

	s.router.HandleFunc("/accounts/{accountId}/inboxes/{inboxId}/rules", s.withRecovery(s.createRule)).Methods("POST")
	s.router.HandleFunc("/accounts/{accountId}/inboxes/{inboxId}/rules", s.withRecovery(s.getRules)).Methods("GET")
	s.router.HandleFunc("/accounts/{accountId}/inboxes/{inboxId}/rules/{ruleId}", s.withRecovery(s.getRule)).Methods("GET")
	s.router.HandleFunc("/accounts/{accountId}/inboxes/{inboxId}/rules/{ruleId}", s.withRecovery(s.updateRule)).Methods("PUT")
	s.router.HandleFunc("/accounts/{accountId}/inboxes/{inboxId}/rules/{ruleId}", s.withRecovery(s.deleteRule)).Methods("DELETE")

	s.router.HandleFunc("/accounts/{accountId}/inboxes/{inboxId}/messages", s.withRecovery(s.getMessages)).Methods("GET")
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.core.Config.HTTPPort, s.router)
}
