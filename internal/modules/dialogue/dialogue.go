package dialogue

import (
	"gorl/internal/logging"
	"gorl/internal/saving"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
	"gointernal.in/yaml.v2"
)

// DialogueNode represents a single piece of dialogue in the game,
// including the text, the emotional profile associated with it, and a flag indicating if it was used.
type DialogueNode struct {
	Text     string           `yaml:"text"`
	Emotions EmotionalProfile `yaml:"emotions"`
}

// EmotionalProfile captures the emotional content of a dialogue statement,
// with each field representing a fraction of the specific emotion present.
type EmotionalProfile struct {
	ImmediateNegative    float32 `yaml:"immediate_negative"`    // Anger, Sadness, Frustration
	AnticipatingNegative float32 `yaml:"anticipating_negative"` // Anxiety, Fear
	ImmediatePositive    float32 `yaml:"immediate_positive"`    // Joy, Excitement, Love
}

// How emotional statement map to emotional responses according to Annes
// personality:
// - Immediate Negative -> Empathy (acknowledge the emotion, but don't try to fix it)
// - Anticipating Negative -> Comfort, Reassurance (provide support)
// - Immediate Positive -> Excitement, Happiness, Encouragement (mirror the emotion)

// Dialogue is a collection of dialogue nodes, with logic to choose the next node to display.
type Dialogue struct {
	dialogueNodes       []*DialogueNode
	currentDialogueNode *DialogueNode
}

func NewDialogueAndLoadData() *Dialogue {
	return &Dialogue{
		dialogueNodes: loadData(),
	}
}

func (d *Dialogue) GetNextNode() *DialogueNode {
	if int(saving.GetInstance().CurrentDialogIndex) >= len(d.dialogueNodes) {
		// eh I think we just gonna loop back to the start
		saving.GetInstance().CurrentDialogIndex = 0
	}
	didx := int(saving.GetInstance().CurrentDialogIndex)
	d.currentDialogueNode = d.dialogueNodes[didx]
	saving.GetInstance().CurrentDialogIndex = int32(didx + 1)
	return d.currentDialogueNode
}

func (d *Dialogue) AttemptEmotionalResponse(empathy, comfort, excitement float32) float32 {
	return rl.Vector3Distance(
		rl.NewVector3(empathy, comfort, excitement),
		rl.NewVector3(
			d.currentDialogueNode.Emotions.ImmediateNegative,
			d.currentDialogueNode.Emotions.AnticipatingNegative,
			d.currentDialogueNode.Emotions.ImmediatePositive))
}

func loadData() []*DialogueNode {
	file, err := os.ReadFile("dialogue/dialogue.yml")
	if err != nil {
		logging.Fatal("Error while loading dialogue.yml file: %v", err)
	}

	var dn []*DialogueNode
	err = yaml.Unmarshal(file, &dn)
	if err != nil {
		logging.Fatal("Error while parsing dialogue.yml file: %v", err)
	}

	return dn
}
