package digest

import (
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-dao/state"
	"github.com/ProtoconNet/mitum-dao/types"
	mitumutil "github.com/ProtoconNet/mitum2/util"
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
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Millisecond*500)
		}
	}
}

func (hd *Handlers) handleDAODesignInGroup(contract string) (interface{}, error) {
	switch design, err := DAOService(hd.database, contract); {
	case err != nil:
		return nil, mitumutil.ErrNotFound.WithMessage(err, "dao service, contract %s", contract)
	case design == nil:
		return nil, mitumutil.ErrNotFound.Errorf("dao service, contract %s", contract)
	default:
		hal, err := hd.buildDAODesignHal(contract, *design)
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

func (hd *Handlers) handleDAOProposal(w http.ResponseWriter, r *http.Request) {
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
		return hd.handleDAOProposalInGroup(contract, proposalID)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Millisecond*500)
		}
	}
}

func (hd *Handlers) handleDAOProposalInGroup(contract, proposalID string) (interface{}, error) {
	switch proposal, err := DAOProposal(hd.database, contract, proposalID); {
	case err != nil:
		return nil, mitumutil.ErrNotFound.WithMessage(err, "proposal, contract %s, proposalID %s", contract, proposalID)
	case proposal == nil:
		return nil, mitumutil.ErrNotFound.Errorf("proposal, contract %s, proposalID %s", contract, proposalID)
	default:
		hal, err := hd.buildDAOProposalHal(contract, proposalID, *proposal)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildDAOProposalHal(contract, proposalID string, proposal state.ProposalStateValue) (currencydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathDAOProposal, "contract", contract, "proposal_id", proposalID)
	if err != nil {
		return nil, err
	}

	hal := currencydigest.NewBaseHal(proposal, currencydigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleDAODelegator(w http.ResponseWriter, r *http.Request) {
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
		return hd.handleDAODelegatorInGroup(contract, proposalID, delegator)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Millisecond*500)
		}
	}
}

func (hd *Handlers) handleDAODelegatorInGroup(contract, proposalID, delegator string) (interface{}, error) {
	switch delegatorInfo, err := DAODelegatorInfo(hd.database, contract, proposalID, delegator); {
	case err != nil:
		return nil, mitumutil.ErrNotFound.WithMessage(err, "delegator info, contract %s, proposalID %s, delegator %s", contract, proposalID, delegator)
	case delegatorInfo == nil:
		return nil, mitumutil.ErrNotFound.Errorf("delegator info, contract %s, proposalID %s, delegator %s", contract, proposalID, delegator)
	default:
		hal, err := hd.buildDAODelegatorHal(contract, proposalID, delegator, *delegatorInfo)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildDAODelegatorHal(
	contract, proposalID, delegator string,
	delegatorInfo types.DelegatorInfo,
) (currencydigest.Hal, error) {
	h, err := hd.combineURL(
		HandlerPathDAODelegator,
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

func (hd *Handlers) handleDAOVoters(w http.ResponseWriter, r *http.Request) {
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
		return hd.handleDAOVotersInGroup(contract, proposalID)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Millisecond*500)
		}
	}
}

func (hd *Handlers) handleDAOVotersInGroup(contract, proposalID string) (interface{}, error) {
	switch voters, err := DAOVoters(hd.database, contract, proposalID); {
	case err != nil:
		return nil, mitumutil.ErrNotFound.WithMessage(err, "voters, contract %s, proposalID %s", contract, proposalID)
	case voters == nil:
		return nil, mitumutil.ErrNotFound.Errorf("voters, contract %s, proposalID %s", contract, proposalID)
	default:
		hal, err := hd.buildDAOVotersHal(contract, proposalID, voters)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildDAOVotersHal(
	contract, proposalID string, voters []types.VoterInfo,
) (currencydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathDAOVoters, "contract", contract, "proposal_id", proposalID)
	if err != nil {
		return nil, err
	}

	hal := currencydigest.NewBaseHal(voters, currencydigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleDAOVotingPowerBox(w http.ResponseWriter, r *http.Request) {
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
		return hd.handleDAOVotingPowerBoxInGroup(contract, proposalID)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Millisecond*500)
		}
	}
}

func (hd *Handlers) handleDAOVotingPowerBoxInGroup(contract, proposalID string) (interface{}, error) {
	switch votingPowerBox, err := DAOVotingPowerBox(hd.database, contract, proposalID); {
	case err != nil:
		return nil, mitumutil.ErrNotFound.WithMessage(err, "voting power box, contract %s, proposalID %s", contract, proposalID)
	case votingPowerBox == nil:
		return nil, mitumutil.ErrNotFound.Errorf("voting power box, contract %s, proposalID %s", contract, proposalID)

	default:
		hal, err := hd.buildDAOVotingPowerBoxHal(contract, proposalID, *votingPowerBox)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildDAOVotingPowerBoxHal(
	contract, proposalID string,
	votingPowerBox types.VotingPowerBox,
) (currencydigest.Hal, error) {
	h, err := hd.combineURL(
		HandlerPathDAOVotingPowerBox,
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
