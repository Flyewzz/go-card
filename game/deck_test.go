package game

import (
	"reflect"
	"sync"
	"testing"
)

func TestNewDeck(t *testing.T) {
	type args struct {
		cards []*Card
	}
	tests := []struct {
		name string
		args args
		want *Deck
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDeck(tt.args.cards); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeck_Push(t *testing.T) {
	tests := []struct {
		name   string
		deck   *Deck
		card   *Card
		result *Deck
	}{
		{
			name: "empty",
			deck: NewDeck([]*Card{}),
			card: &Card{ID: 1},
			result: &Deck{
				Cards: []*Card{
					{ID: 1},
				},
			},
		},
		{
			name: "one number",
			deck: NewDeck([]*Card{
				{ID: 941},
			}),
			card: &Card{ID: 43},
			result: &Deck{
				Cards: []*Card{
					{ID: 43},
					{ID: 941},
				},
			},
		},
		{
			name: "tree number",
			deck: NewDeck([]*Card{
				{ID: 32},
				{ID: 1},
				{ID: 16},
			}),
			card: &Card{ID: 9},
			result: &Deck{
				Cards: []*Card{
					{ID: 9},
					{ID: 32},
					{ID: 1},
					{ID: 16},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := tt.deck
			this.Push(tt.card)
			if !reflect.DeepEqual(tt.deck.Cards, tt.result.Cards) {
				t.Errorf("Expected [%s], got [%s]", tt.result.Print(), tt.deck.Print())
			}
		})
	}
}

func generateCards(low, high int) []*Card {
	var cards []*Card
	for v := low; v <= high; v++ {
		cards = append(cards, &Card{
			ID: v,
		})
	}
	return cards
}
func TestDeck_GetCard(t *testing.T) {
	type args struct {
		position int
	}
	tests := []struct {
		name    string
		deck    *Deck
		args    args
		want    *Card
		wantErr bool
	}{
		{
			name: "First card",
			deck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
				{ID: 10},
			}),
			args: args{
				position: 0,
			},
			want:    &Card{ID: 34},
			wantErr: false,
		},
		{
			name: "Last card",
			deck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
				{ID: 9},
			}),
			args: args{
				position: 3,
			},
			want:    &Card{ID: 9},
			wantErr: false,
		},
		{
			name: "Middle card",
			deck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
				{ID: 10},
				{ID: 18},
				{ID: 67},
			}),
			args: args{
				position: 3,
			},
			want:    &Card{ID: 10},
			wantErr: false,
		},
		{
			name: "Negative position",
			deck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
			}),
			args: args{
				position: -1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Too big position",
			deck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
				{ID: 10},
				{ID: 18},
				{ID: 67},
			}),
			args: args{
				position: 15,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := tt.deck
			got, err := this.GetCard(tt.args.position)
			if (err != nil) != tt.wantErr {
				t.Errorf("Deck.GetCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Deck.GetCard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeck_Pop(t *testing.T) {
	type args struct {
		position int
	}
	tests := []struct {
		name       string
		deck       *Deck
		resultDeck *Deck
		args       args
		want       *Card
		wantErr    bool
	}{
		{
			name: "First card",
			deck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
				{ID: 10},
			}),
			resultDeck: NewDeck([]*Card{
				{ID: 1},
				{ID: 12},
				{ID: 10},
			}),
			args: args{
				position: 0,
			},
			want:    &Card{ID: 34},
			wantErr: false,
		},
		{
			name: "Last card",
			deck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
				{ID: 9},
			}),
			resultDeck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
			}),
			args: args{
				position: 3,
			},
			want:    &Card{ID: 9},
			wantErr: false,
		},
		{
			name: "Middle card",
			deck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
				{ID: 10},
				{ID: 18},
				{ID: 67},
			}),
			resultDeck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
				{ID: 18},
				{ID: 67},
			}),
			args: args{
				position: 3,
			},
			want:    &Card{ID: 10},
			wantErr: false,
		},
		{
			name: "Negative position",
			deck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
			}),
			resultDeck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
			}),
			args: args{
				position: -1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Too big position",
			deck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
				{ID: 10},
				{ID: 18},
				{ID: 67},
			}),
			resultDeck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
				{ID: 10},
				{ID: 18},
				{ID: 67},
			}),
			args: args{
				position: 15,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := tt.deck
			got, err := this.Pop(tt.args.position)
			if (err != nil) != tt.wantErr {
				t.Errorf("Deck.GetCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Deck.GetCard() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(tt.deck.Cards, tt.resultDeck.Cards) {
				t.Errorf("Expected %s, but got %s",
					tt.deck.Print(), tt.resultDeck.Print())
			}
		})
	}
}

func TestDeck_Shuffle(t *testing.T) {
	type fields struct {
		mtx   *sync.Mutex
		Cards []*Card
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Deck{
				mtx:   tt.fields.mtx,
				Cards: tt.fields.Cards,
			}
			this.Shuffle()
		})
	}
}

func TestDeck_MoveCard(t *testing.T) {
	type args struct {
		position int
	}
	tests := []struct {
		name        string
		deck        *Deck
		anotherDeck *Deck
		args        args
		want        *Deck
		wantErr     bool
	}{
		{
			name: "First card",
			deck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
				{ID: 10},
			}),
			anotherDeck: NewDeck([]*Card{
				{ID: 7},
				{ID: 3},
				{ID: 0},
			}),
			args: args{
				position: 0,
			},
			want: NewDeck([]*Card{
				{ID: 34},
				{ID: 7},
				{ID: 3},
				{ID: 0},
			}),
			wantErr: false,
		},
		{
			name: "Last card",
			deck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
				{ID: 9},
			}),
			anotherDeck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
			}),
			args: args{
				position: 3,
			},
			want: NewDeck([]*Card{
				{ID: 9},
				{ID: 34},
				{ID: 1},
				{ID: 12},
			}),
			wantErr: false,
		},
		{
			name: "Middle card",
			deck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
				{ID: 10},
				{ID: 18},
				{ID: 67},
			}),
			anotherDeck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
				{ID: 18},
				{ID: 67},
			}),
			args: args{
				position: 3,
			},
			want: NewDeck([]*Card{
				{ID: 10},
				{ID: 34},
				{ID: 1},
				{ID: 12},
				{ID: 18},
				{ID: 67},
			}),
			wantErr: false,
		},
		{
			name: "Negative position",
			deck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
			}),
			anotherDeck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
			}),
			args: args{
				position: -1,
			},
			want: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
			}),
			wantErr: true,
		},
		{
			name: "Too big position",
			deck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
				{ID: 10},
				{ID: 18},
				{ID: 67},
			}),
			anotherDeck: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
				{ID: 10},
				{ID: 18},
				{ID: 67},
			}),
			args: args{
				position: 15,
			},
			want: NewDeck([]*Card{
				{ID: 34},
				{ID: 1},
				{ID: 12},
				{ID: 10},
				{ID: 18},
				{ID: 67},
			}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := tt.deck
			if err := this.MoveCard(tt.args.position, tt.anotherDeck); (err != nil) != tt.wantErr {
				t.Errorf("Deck.MoveCard() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.anotherDeck.Cards, tt.want.Cards) {
				t.Errorf("Expected %s, but got %s",
					tt.anotherDeck.Print(), tt.want.Print())
			}
		})
	}
}

func TestDeck_Print(t *testing.T) {
	type fields struct {
		mtx   *sync.Mutex
		Cards []*Card
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Deck{
				mtx:   tt.fields.mtx,
				Cards: tt.fields.Cards,
			}
			if got := this.Print(); got != tt.want {
				t.Errorf("Deck.Print() = %v, want %v", got, tt.want)
			}
		})
	}
}
