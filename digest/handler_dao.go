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

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleDAODesignInGroup(contract)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleDAODesignInGroup(contract string) (interface{}, error) {
	switch design, err := DAOService(hd.database, contract); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildDAODesignHal(contract, design)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildDAODesignHal(contract string, design types.Design) (currencydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathDAOService, "contract", contract)
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

	proposalID, err, status := parseRequest(w, r, "proposal_id")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleProposalInGroup(contract, proposalID)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleProposalInGroup(contract, proposalID string) (interface{}, error) {
	switch proposal, err := Proposal(hd.database, contract, proposalID); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildProposalHal(contract, proposalID, proposal)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildProposalHal(contract, proposalID string, proposal state.ProposalStateValue) (currencydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathDAOService, "contract", contract, "proposal_id", proposalID)
	if err != nil {
		return nil, err
	}

	hal := currencydigest.NewBaseHal(proposal, currencydigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleDelegator(w http.ResponseWriter, r *http.Request) {
	cacheKey := currencydigest.CacheKeyPath(r)
	if err := currencydigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	proposalID, err, status := parseRequest(w, r, "proposal_id")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	delegator, err, status := parseRequest(w, r, "address")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleDelegatorInGroup(contract, proposalID, delegator)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleDelegatorInGroup(contract, proposalID, delegator string) (interface{}, error) {
	switch delegatorInfo, err := DelegatorInfo(hd.database, contract, proposalID, delegator); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildDelegatorHal(contract, proposalID, delegator, delegatorInfo)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildDelegatorHal(
	contract, proposalID, delegator string,
	delegatorInfo types.DelegatorInfo,
) (currencydigest.Hal, error) {
	h, err := hd.combineURL(
		HandlerPathDelegator,
		"contract", contract,
		"proposal_id", proposalID,
		"address", delegator,
	)
	if err != nil {
		return nil, err
	}

	hal := currencydigest.NewBaseHal(delegatorInfo, currencydigest.NewHalLink(h, nil))

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

	proposalID, err, status := parseRequest(w, r, "proposal_id")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleVotersInGroup(contract, proposalID)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleVotersInGroup(contract, proposalID string) (interface{}, error) {
	switch voters, err := Voters(hd.database, contract, proposalID); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildVotersHal(contract, proposalID, voters)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildVotersHal(
	contract, proposalID string, voters []types.VoterInfo,
) (currencydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathVoters, "contract", contract, "proposal_id", proposalID)
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

	proposalID, err, status := parseRequest(w, r, "proposal_id")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleVotingPowerBoxInGroup(contract, proposalID)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleVotingPowerBoxInGroup(contract, proposalID string) (interface{}, error) {
	switch votingPowerBox, err := VotingPowerBox(hd.database, contract, proposalID); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildVotingPowerBoxHal(contract, proposalID, votingPowerBox)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildVotingPowerBoxHal(
	contract, proposalID string,
	votingPowerBox types.VotingPowerBox,
) (currencydigest.Hal, error) {
	h, err := hd.combineURL(
		HandlerPathVotingPowerBox,
		"contract", contract,
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
