package integration_test

import (
	"testing"

	boardgamegeek "github.com/Boughtmanatee5/BoardGameGeekClient"
	"github.com/stretchr/testify/assert"
)

// Game: Catan family
const catanFamilyId int64 = 3

// Crowdfunding: Kickstarter family
const kickstarterFamilyId int64 = 8374

func TestGetFamilyItemsById(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.FamilyItemRequest{
		IDs: []int64{catanFamilyId},
	}

	response, err := client.GetFamilyItems(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, response.Items, 1)
	assert.Equal(t, catanFamilyId, response.Items[0].Id)
	assert.Equal(t, string(boardgamegeek.FamilyItemTypeBoardGameFamily), response.Items[0].Type)
}

func TestGetFamilyItemsPopulatesFields(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.FamilyItemRequest{
		IDs: []int64{catanFamilyId},
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
	var linkedIds []int64
	for _, link := range item.Links {
		linkedIds = append(linkedIds, link.Id)
	}
	assert.Contains(t, linkedIds, catanId)
	assert.True(t, item.Links[0].Inbound)
}

func TestGetFamilyItemsByMultipleIds(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.FamilyItemRequest{
		IDs: []int64{catanFamilyId, kickstarterFamilyId},
	}

	response, err := client.GetFamilyItems(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, response.Items, 2)

	var returnedIds []int64
	for _, item := range response.Items {
		returnedIds = append(returnedIds, item.Id)
	}
	assert.ElementsMatch(t, []int64{catanFamilyId, kickstarterFamilyId}, returnedIds)
}

func TestGetFamilyItemsByType(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.FamilyItemRequest{
		IDs:   []int64{catanFamilyId},
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
	request := &boardgamegeek.FamilyItemRequest{
		IDs:   []int64{catanFamilyId},
		Types: []boardgamegeek.FamilyItemType{boardgamegeek.FamilyItemTypeRPG},
	}

	response, err := client.GetFamilyItems(ctx, request)
	assert.NoError(t, err)
	assert.Empty(t, response.Items)
}

func TestGetFamilyItemsByInvalidId(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.FamilyItemRequest{
		IDs: []int64{999999999},
	}

	response, err := client.GetFamilyItems(ctx, request)
	assert.NoError(t, err)
	assert.Empty(t, response.Items)
}

func TestGetFamilyItemsByMixOfValidAndInvalidIds(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.FamilyItemRequest{
		IDs: []int64{catanFamilyId, 999999999},
	}

	response, err := client.GetFamilyItems(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, response.Items, 1)
	assert.Equal(t, catanFamilyId, response.Items[0].Id)
}
