package handlers

import (
    "database/sql"
    "net/http"
    "time"

    "voting-system/db"
    "voting-system/models"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

// CreatePoll handles POST /api/polls
func CreatePoll(c *gin.Context) {
    var req models.CreatePollRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Start transaction
    tx, err := db.DB.Begin()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create poll"})
        return
    }
    defer tx.Rollback()

    // Insert poll
    pollID := uuid.New()
    _, err = tx.Exec(`
        INSERT INTO polls (id, question, expires_at)
        VALUES ($1, $2, $3)`,
        pollID, req.Question, req.ExpiresAt)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create poll"})
        return
    }

    // Insert options
    for _, optValue := range req.Options {
        optID := uuid.New()
        _, err = tx.Exec(`
            INSERT INTO options (id, poll_id, value)
            VALUES ($1, $2, $3)`,
            optID, pollID, optValue)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add options"})
            return
        }
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save poll"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "poll_id": pollID,
        "message": "Poll created successfully",
    })
}

// GetPoll handles GET /api/polls/:id
func GetPoll(c *gin.Context) {
    pollID := c.Param("id")
    
    // Get poll details
    var poll models.Poll
    err := db.DB.QueryRow(`
        SELECT id, question, created_at, expires_at, is_active
        FROM polls WHERE id = $1`,
        pollID).Scan(&poll.ID, &poll.Question, &poll.CreatedAt, &poll.ExpiresAt, &poll.IsActive)
    
    if err == sql.ErrNoRows {
        c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
        return
    }
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }

    // Get options with vote counts
    rows, err := db.DB.Query(`
        SELECT id, value, votes_count
        FROM options WHERE poll_id = $1`,
        pollID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch options"})
        return
    }
    defer rows.Close()

    var options []models.Option
    for rows.Next() {
        var opt models.Option
        err := rows.Scan(&opt.ID, &opt.Value, &opt.VotesCount)
        if err != nil {
            continue
        }
        options = append(options, opt)
    }

    c.JSON(http.StatusOK, gin.H{
        "poll":    poll,
        "options": options,
    })
}

// Vote handles POST /api/vote
func Vote(c *gin.Context) {
    var vote models.Vote
    if err := c.ShouldBindJSON(&vote); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Get voter's IP address
    ipAddress := c.ClientIP()
    vote.IPAddress = ipAddress

    // Check if already voted
    var exists bool
    err := db.DB.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM votes 
            WHERE poll_id = $1 AND ip_address = $2
        )`, vote.PollID, vote.IPAddress).Scan(&exists)
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check vote status"})
        return
    }

    if exists {
        c.JSON(http.StatusConflict, gin.H{"error": "You have already voted in this poll"})
        return
    }

    // Start transaction for vote
    tx, err := db.DB.Begin()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process vote"})
        return
    }
    defer tx.Rollback()

    // Insert vote
    _, err = tx.Exec(`
        INSERT INTO votes (poll_id, option_id, ip_address)
        VALUES ($1, $2, $3)`,
        vote.PollID, vote.OptionID, vote.IPAddress)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record vote"})
        return
    }

    // Increment vote count for option
    _, err = tx.Exec(`
        UPDATE options 
        SET votes_count = votes_count + 1 
        WHERE id = $1`,
        vote.OptionID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vote count"})
        return
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save vote"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Vote recorded successfully"})
}

// GetResults handles GET /api/results/:pollId
func GetResults(c *gin.Context) {
    pollID := c.Param("pollId")
    
    rows, err := db.DB.Query(`
        SELECT o.value, o.votes_count,
               ROUND(100.0 * o.votes_count / NULLIF(SUM(o.votes_count) OVER(), 0), 2) as percentage
        FROM options o
        WHERE o.poll_id = $1
        ORDER BY o.votes_count DESC`,
        pollID)
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch results"})
        return
    }
    defer rows.Close()

    var results []gin.H
    for rows.Next() {
        var option string
        var votes int
        var percentage float64
        rows.Scan(&option, &votes, &percentage)
        
        results = append(results, gin.H{
            "option":     option,
            "votes":      votes,
            "percentage": percentage,
        })
    }

    c.JSON(http.StatusOK, gin.H{"results": results})
}