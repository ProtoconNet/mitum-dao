package digest

import (
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-dao/state"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func (hd *Handlers) handleDAOService(w http.ResponseWriter, r *http.Request) {
	cacheKey := currencydigest.CacheKeyPath(r)
	if err := currencydigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	service, err, status := parseRequest(w, r, "dao_id")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleDAODesignInGroup(contract, service)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleDAODesignInGroup(contract, service string) (interface{}, error) {
	switch design, err := DAOService(hd.database, contract, service); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildDAODesignHal(contract, service, design)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildDAODesignHal(contract, service string, design types.Design) (currencydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathDAOService, "contract", contract, "dao_id", service)
	if err != nil {
		return nil, err
	}

	hal := currencydigest.NewBaseHal(design, currencydigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleProposal(w http.ResponseWriter, r *http.Request) {
	cacheKey := currencydigest.CacheKeyPath(r)
	if err := currencydigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	daoID, err, status := parseRequest(w, r, "dao_id")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	proposalID, err, status := parseRequest(w, r, "proposal_id")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleProposalInGroup(contract, daoID, proposalID)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleProposalInGroup(contract, daoID, proposalID string) (interface{}, error) {
	switch proposal, err := Proposal(hd.database, contract, daoID, proposalID); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildProposalHal(contract, daoID, proposalID, proposal)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildProposalHal(contract, daoID, proposalID string, proposal state.ProposalStateValue) (currencydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathDAOService, "contract", contract, "dao_id", daoID, "proposal_id", proposalID)
	if err != nil {
		return nil, err
	}

	hal := currencydigest.NewBaseHal(proposal, currencydigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleDelegators(w http.ResponseWriter, r *http.Request) {
	cacheKey := currencydigest.CacheKeyPath(r)
	if err := currencydigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	serviceID, err, status := parseRequest(w, r, "dao_id")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	proposalID, err, status := parseRequest(w, r, "proposal_id")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleDelegatorsInGroup(contract, serviceID, proposalID)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleDelegatorsInGroup(contract, serviceID, proposalID string) (interface{}, error) {
	switch delegators, err := Delegators(hd.database, contract, serviceID, proposalID); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildDelegatorsHal(contract, serviceID, proposalID, delegators)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildDelegatorsHal(
	contract, serviceID, proposalID string,
	delegators []types.DelegatorInfo,
) (currencydigest.Hal, error) {
	h, err := hd.combineURL(
		HandlerPathDelegators,
		"contract", contract,
		"dao_id", serviceID,
		"proposal_id", proposalID,
	)
	if err != nil {
		return nil, err
	}

	hal := currencydigest.NewBaseHal(delegators, currencydigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleVoters(w http.ResponseWriter, r *http.Request) {
	cacheKey := currencydigest.CacheKeyPath(r)
	if err := currencydigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	daoID, err, status := parseRequest(w, r, "dao_id")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	proposalID, err, status := parseRequest(w, r, "proposal_id")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleVotersInGroup(contract, daoID, proposalID)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleVotersInGroup(contract, daoID, proposalID string) (interface{}, error) {
	switch voters, err := Voters(hd.database, contract, daoID, proposalID); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildVotersHal(contract, daoID, proposalID, voters)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildVotersHal(
	contract, daoID, proposalID string, voters []types.VoterInfo,
) (currencydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathVoters, "contract", contract, "dao_id", daoID, "proposal_id", proposalID)
	if err != nil {
		return nil, err
	}

	hal := currencydigest.NewBaseHal(voters, currencydigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleVotingPowerBox(w http.ResponseWriter, r *http.Request) {
	cacheKey := currencydigest.CacheKeyPath(r)
	if err := currencydigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	daoID, err, status := parseRequest(w, r, "dao_id")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	proposalID, err, status := parseRequest(w, r, "proposal_id")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleVotingPowerBoxInGroup(contract, daoID, proposalID)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleVotingPowerBoxInGroup(contract, daoID, proposalID string) (interface{}, error) {
	switch votingPowerBox, err := VotingPowerBox(hd.database, contract, daoID, proposalID); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildVotingPowerBoxHal(contract, daoID, proposalID, votingPowerBox)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildVotingPowerBoxHal(
	contract, daoID, proposalID string,
	votingPowerBox types.VotingPowerBox,
) (currencydigest.Hal, error) {
	h, err := hd.combineURL(
		HandlerPathVotingPowerBox,
		"contract", contract,
		"dao_id", daoID,
		"proposal_id", proposalID,
	)
	if err != nil {
		return nil, err
	}

	hal := currencydigest.NewBaseHal(votingPowerBox, currencydigest.NewHalLink(h, nil))

	return hal, nil
}

func parseRequest(_ http.ResponseWriter, r *http.Request, v string) (string, error, int) {
	s, found := mux.Vars(r)[v]
	if !found {
		return "", errors.Errorf("empty %s", v), http.StatusNotFound
	}

	s = strings.TrimSpace(s)
	if len(s) < 1 {
		return "", errors.Errorf("empty %s", v), http.StatusBadRequest
	}
	return s, nil, http.StatusOK
}
