package integration_test

import (
	"testing"

	"github.com/Boughtmanatee5/go-boardgamegeek"
	"github.com/stretchr/testify/assert"
)

// Game: Catan family
const catanFamilyID int64 = 3

// Crowdfunding: Kickstarter family
const kickstarterFamilyID int64 = 8374

func TestGetFamilyItemsByID(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.FamilyItemsOptions{
		IDs: []int64{catanFamilyID},
	}

	response, err := client.GetFamilyItems(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, response.Items, 1)
	assert.Equal(t, catanFamilyID, response.Items[0].ID)
	assert.Equal(t, string(boardgamegeek.FamilyItemTypeBoardGameFamily), response.Items[0].Type)
}

func TestGetFamilyItemsPopulatesFields(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.FamilyItemsOptions{
		IDs: []int64{catanFamilyID},
	}

	response, err := client.GetFamilyItems(ctx, request)
	assert.NoError(t, err)

	item := response.Items[0]
	assert.NotEmpty(t, item.Thumbnail)
	assert.NotEmpty(t, item.Image)
	assert.NotEmpty(t, item.Description)
	assert.NotEmpty(t, item.Names)
	assert.Equal(t, "primary", item.Names[0].Type)
	assert.Equal(t, "Game: Catan", item.Names[0].Value)
	assert.NotEmpty(t, item.Links)

	// Catan (13) is a member of the Game: Catan family
	var linkedIDs []int64
	for _, link := range item.Links {
		linkedIDs = append(linkedIDs, link.ID)
	}
	assert.Contains(t, linkedIDs, catanID)
	assert.True(t, item.Links[0].Inbound)
}

func TestGetFamilyItemsByMultipleIDs(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.FamilyItemsOptions{
		IDs: []int64{catanFamilyID, kickstarterFamilyID},
	}

	response, err := client.GetFamilyItems(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, response.Items, 2)

	var returnedIDs []int64
	for _, item := range response.Items {
		returnedIDs = append(returnedIDs, item.ID)
	}
	assert.ElementsMatch(t, []int64{catanFamilyID, kickstarterFamilyID}, returnedIDs)
}

func TestGetFamilyItemsByType(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.FamilyItemsOptions{
		IDs:   []int64{catanFamilyID},
		Types: []boardgamegeek.FamilyItemType{boardgamegeek.FamilyItemTypeBoardGameFamily},
	}

	response, err := client.GetFamilyItems(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, response.Items, 1)
	assert.Equal(t, string(boardgamegeek.FamilyItemTypeBoardGameFamily), response.Items[0].Type)
}

// Edge cases

func TestGetFamilyItemsByNonMatchingType(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.FamilyItemsOptions{
		IDs:   []int64{catanFamilyID},
		Types: []boardgamegeek.FamilyItemType{boardgamegeek.FamilyItemTypeRPG},
	}

	response, err := client.GetFamilyItems(ctx, request)
	assert.NoError(t, err)
	assert.Empty(t, response.Items)
}

func TestGetFamilyItemsByInvalidID(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.FamilyItemsOptions{
		IDs: []int64{999999999},
	}

	response, err := client.GetFamilyItems(ctx, request)
	assert.NoError(t, err)
	assert.Empty(t, response.Items)
}

func TestGetFamilyItemsByMixOfValidAndInvalidIDs(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.FamilyItemsOptions{
		IDs: []int64{catanFamilyID, 999999999},
	}

	response, err := client.GetFamilyItems(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, response.Items, 1)
	assert.Equal(t, catanFamilyID, response.Items[0].ID)
}
