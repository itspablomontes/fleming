package audit

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleVerifyRoot(c *gin.Context) {
	if h.chainClient == nil {
		c.JSON(http.StatusNotImplemented, gin.H{
			"error": "chain anchoring is not configured (set ANCHOR_RPC_URL, ANCHOR_CONTRACT_ADDRESS, ANCHOR_PRIVATE_KEY)",
		})
		return
	}

	address, exists := c.Get("user_address")
	actor, ok := address.(string)
	if !exists || !ok || actor == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	root := strings.TrimSpace(c.Param("root"))
	if root == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "root is required"})
		return
	}

	timestamp, err := h.chainClient.VerifyRoot(c.Request.Context(), root)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to verify root"})
		return
	}

	anchored := timestamp != 0

	var txHash *string
	var blockNumber *uint64

	// Prefer DB (actor-scoped) data if we have a matching batch.
	batch, err := h.service.GetBatchByRoot(c.Request.Context(), actor, root)
	if err == nil && batch != nil && batch.AnchorTxHash != nil && batch.AnchorBlockNumber != nil {
		txHash = batch.AnchorTxHash
		blockNumber = batch.AnchorBlockNumber
	} else if anchored {
		// Otherwise fall back to event lookup (useful when the user is verifying a root
		// they didn't anchor via this backend instance).
		ev, found, err := h.chainClient.FindRootAnchoredEvent(c.Request.Context(), root)
		if err == nil && found && ev != nil {
			txHash = &ev.TxHash
			blockNumber = &ev.BlockNumber
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"anchored":    anchored,
		"timestamp":   timestamp,
		"blockNumber": blockNumber,
		"txHash":      txHash,
	})
}
