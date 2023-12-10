package choose

import (
	"context"
	"reflect"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/uuid"
)

var ChooseMap map[string]struct {
	Type     reflect.Type
	Retrieve func(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) ([]uuid.UUID, error)
}

func init() {
	ChooseMap = map[string]struct {
		Type     reflect.Type
		Retrieve func(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) ([]uuid.UUID, error)
	}{
		AdjacentChoice: {
			Type:     reflect.TypeOf(&AdjacentArgs{}),
			Retrieve: RetrieveAdjacent,
		},
		CodexChoice: {
			Type:     reflect.TypeOf(&CodexArgs{}),
			Retrieve: RetrieveCodex,
		},
		CompositeChoice: {
			Type:     reflect.TypeOf(&CompositeArgs{}),
			Retrieve: RetrieveComposite,
		},
		ConnectedChoice: {
			Type:     reflect.TypeOf(&ConnectedArgs{}),
			Retrieve: RetrieveConnected,
		},
		CurrentPlayerChoice: {
			Type:     reflect.TypeOf(&CurrentPlayerArgs{}),
			Retrieve: RetrieveCurrentPlayer,
		},
		HookEventUnitChoice: {
			Type:     reflect.TypeOf(&HookEventUnitArgs{}),
			Retrieve: RetrieveHookEventUnit,
		},
		OpposingPlayerChoice: {
			Type:     reflect.TypeOf(&OpposingPlayerArgs{}),
			Retrieve: RetrieveOpposingPlayer,
		},
		OwnedTilesChoice: {
			Type:     reflect.TypeOf(&OwnedTilesArgs{}),
			Retrieve: RetrieveOwnedTiles,
		},
		OwnedUnitsChoice: {
			Type:     reflect.TypeOf(&OwnedUnitsArgs{}),
			Retrieve: RetrieveOwnedUnits,
		},
		RandomChoice: {
			Type:     reflect.TypeOf(&RandomArgs{}),
			Retrieve: RetrieveRandom,
		},
		SelfOwnerChoice: {
			Type:     reflect.TypeOf(&SelfOwnerArgs{}),
			Retrieve: RetrieveSelfOwner,
		},
		SelfChoice: {
			Type:     reflect.TypeOf(&SelfArgs{}),
			Retrieve: RetrieveSelf,
		},
		TargetOwnerChoice: {
			Type:     reflect.TypeOf(&TargetOwnerArgs{}),
			Retrieve: RetrieveTargetOwner,
		},
		TargetChoice: {
			Type:     reflect.TypeOf(&TargetArgs{}),
			Retrieve: RetrieveTarget,
		},
		TilesChoice: {
			Type:     reflect.TypeOf(&TilesArgs{}),
			Retrieve: RetrieveTiles,
		},
		UnitsChoice: {
			Type:     reflect.TypeOf(&UnitsArgs{}),
			Retrieve: RetrieveUnits,
		},
		UUIDChoice: {
			Type:     reflect.TypeOf(&UUIDArgs{}),
			Retrieve: RetrieveUUID,
		},
	}
}
