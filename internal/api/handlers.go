package api

import (
	"encoding/json"
	"mercury/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) createAccount(w http.ResponseWriter, r *http.Request) {
	var account models.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		s.core.HandleError(w, err, http.StatusBadRequest)
		return
	}

	if err := s.core.AccountService.Create(&account); err != nil {
		s.core.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

func (s *Server) getAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := s.core.AccountService.List()
	if err != nil {
		s.core.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(accounts)
}

func (s *Server) getAccount(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	account, err := s.core.AccountService.Get(id)
	if err != nil {
		s.core.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	if account == nil {
		s.core.HandleError(w, err, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(account)
}

func (s *Server) updateAccount(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var account models.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		s.core.HandleError(w, err, http.StatusBadRequest)
		return
	}
	account.ID = id

	if err := s.core.AccountService.Update(&account); err != nil {
		s.core.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) deleteAccount(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := s.core.AccountService.Delete(id); err != nil {
		s.core.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) createInbox(w http.ResponseWriter, r *http.Request) {
	accountID, _ := strconv.Atoi(mux.Vars(r)["accountId"])
	var inbox models.Inbox
	if err := json.NewDecoder(r.Body).Decode(&inbox); err != nil {
		s.core.HandleError(w, err, http.StatusBadRequest)
		return
	}
	inbox.AccountID = accountID

	if err := s.core.InboxService.Create(&inbox); err != nil {
		s.core.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(inbox)
}

func (s *Server) getInboxes(w http.ResponseWriter, r *http.Request) {
	accountID, _ := strconv.Atoi(mux.Vars(r)["accountId"])
	inboxes, err := s.core.InboxService.GetByAccountID(accountID)
	if err != nil {
		s.core.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(inboxes)
}

func (s *Server) getInbox(w http.ResponseWriter, r *http.Request) {
	inboxID, _ := strconv.Atoi(mux.Vars(r)["inboxId"])
	inbox, err := s.core.InboxService.Get(inboxID)
	if err != nil {
		s.core.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	if inbox == nil {
		s.core.HandleError(w, err, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(inbox)
}

func (s *Server) updateInbox(w http.ResponseWriter, r *http.Request) {
	inboxID, _ := strconv.Atoi(mux.Vars(r)["inboxId"])
	var inbox models.Inbox
	if err := json.NewDecoder(r.Body).Decode(&inbox); err != nil {
		s.core.HandleError(w, err, http.StatusBadRequest)
		return
	}
	inbox.ID = inboxID

	if err := s.core.InboxService.Update(&inbox); err != nil {
		s.core.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) deleteInbox(w http.ResponseWriter, r *http.Request) {
	inboxID, _ := strconv.Atoi(mux.Vars(r)["inboxId"])
	if err := s.core.InboxService.Delete(inboxID); err != nil {
		s.core.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) createRule(w http.ResponseWriter, r *http.Request) {
	inboxID, _ := strconv.Atoi(mux.Vars(r)["inboxId"])
	var rule models.Rule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		s.core.HandleError(w, err, http.StatusBadRequest)
		return
	}
	rule.InboxID = inboxID

	if err := s.core.RuleService.Create(&rule); err != nil {
		s.core.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rule)
}

func (s *Server) getRules(w http.ResponseWriter, r *http.Request) {
	inboxID, _ := strconv.Atoi(mux.Vars(r)["inboxId"])
	rules, err := s.core.RuleService.GetByInboxID(inboxID)
	if err != nil {
		s.core.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rules)
}

func (s *Server) getRule(w http.ResponseWriter, r *http.Request) {
	ruleID, _ := strconv.Atoi(mux.Vars(r)["ruleId"])
	rule, err := s.core.RuleService.Get(ruleID)
	if err != nil {
		s.core.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	if rule == nil {
		s.core.HandleError(w, err, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(rule)
}

func (s *Server) updateRule(w http.ResponseWriter, r *http.Request) {
	ruleID, _ := strconv.Atoi(mux.Vars(r)["ruleId"])
	var rule models.Rule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		s.core.HandleError(w, err, http.StatusBadRequest)
		return
	}
	rule.ID = ruleID

	if err := s.core.RuleService.Update(&rule); err != nil {
		s.core.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) deleteRule(w http.ResponseWriter, r *http.Request) {
	ruleID, _ := strconv.Atoi(mux.Vars(r)["ruleId"])
	if err := s.core.RuleService.Delete(ruleID); err != nil {
		s.core.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) getMessages(w http.ResponseWriter, r *http.Request) {
	inboxID, _ := strconv.Atoi(mux.Vars(r)["inboxId"])
	messages, err := s.core.MessageService.GetByInboxID(inboxID)
	if err != nil {
		s.core.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(messages)
}
