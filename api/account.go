package api

import (
	"database/sql"
	"net/http"

	db "github.com/Richd0tcom/SafeX-Pay/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
)

type createAccountRequest struct {
	Owner string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

type getAccountRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type listAccountsRequest struct {
	Page int32 `form:"page" binding:"required,min=1"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest

	err:= ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, buildErrorResponse(err))
		return
	}

	arg:= db.CreateAccountParams{
		Owner: req.Owner,
		Balance: 0,
		Currency: req.Currency,
	}

	account, err:= server.store.CreateAccount(ctx,arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, buildErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, account)
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest

	err:= ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, buildErrorResponse(err))
		return
	}

	parsedId, _:= uuid.Parse(req.ID)


	account, err:= server.store.GetAccount(ctx,parsedId)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, buildErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, buildErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}


func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest

	err:= ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, buildErrorResponse(err))
		return
	}

	args:=db.ListAccountsParams{
		Limit: 10,
		Offset: (req.Page -1)*10 ,
	}

	accounts, err:= server.store.ListAccounts(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, buildErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}
