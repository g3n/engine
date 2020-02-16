package gui

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/window"
	"time"
)

// Tooltip states
const (
	tooltipOff = iota
	tooltipPending
	tooltipActive
)

// Tooltip positioning
type tooltipAlign int

// Tooltip location identifiers
const (
	tooltipOffset = tooltipAlign(iota)
	tooltipFixed
	tooltipCustom
)

// The ID for the subscriptions
const tooltipID = "tooltip"

type Tooltip struct {
	// relations
	panel IPanel // This IPanel is the actual visual tooltip

	// settings
	delay  time.Duration // the delay duration
	follow bool          // Tooltip follows cursor

	align   tooltipAlign                                 // location identifier
	custom  func(*window.CursorEvent) (float32, float32) // for the identifier tooltipCustom
	offsetX float32                                      // for the identifier tooltipOffset
	offsetY float32                                      //  "	"		"			"

	// internals
	lastevent  *window.CursorEvent // the last event received from OnCursor
	subscribed bool                // if we are subscribed to OnCursor event for the Window().win
	status     int                 // current status of tooltip
	cancel     chan bool           // to cancel a tooltip during the tooltipPending state
}

// Creates a tooltip, the specified IPanel is what will pop up. By default the tooltip will be visible after one second
// and placed around the cursor.
func NewTooltip(panel IPanel) *Tooltip {

	t := new(Tooltip)
	t.panel = panel
	t.delay = time.Second
	t.cancel = make(chan bool)
	t.offsetX = 16
	t.offsetY = 16
	return t
}

// Sets the amount of time that has to be passed before drawing the tooltip.
// If the argument is <= 0, the tooltip will be drawn immediately upon hovering the target.
// The default delay is one second.
func (t *Tooltip) SetDelay(milliseconds int) {
	if milliseconds <= 0 {
		t.delay = 0
	}
	t.delay = time.Millisecond * time.Duration(milliseconds)

}

// Sets the offset between the cursor location and the tooltip.
func (t *Tooltip) SetPositionOffset(x, y float32) {
	t.align = tooltipOffset
	t.offsetX = x
	t.offsetY = y
}

// The tooltip will be fixed to this position.
func (t *Tooltip) SetPositionFixed(x, y float32) {
	t.align = tooltipFixed
	t.panel.SetPosition(x, y)
}

// Registers a custom positioning function for the tooltip.
// The returned values will define the position of the tooltip.
func (t *Tooltip) SetPositionCustom(fn func(event *window.CursorEvent) (newX float32, newY float32)) {

	t.align = tooltipCustom
	t.custom = fn
}

// If set to true, the tooltip will follow the cursor. On false, the location of the tooltip will be based
// upon the last cursor event during the delay state. Has no effect when SetPositionFixed() is used.
func (t *Tooltip) SetFollow(state bool) {
	t.follow = state
}

// this function is called by (IPanel).SetTooltip(...), should not be exposed or called by any other functions
func (t *Tooltip) assign(forPanel IPanel) {

	// When the cursor enters the Panel
	forPanel.SubscribeID(OnCursorEnter, tooltipID, func(s string, i interface{}) {
		if t.subscribed {
			return
		}

		// We let the Manager.win subscription update the OnCursor movement
		// Why? Because if a nested child of forPanel is subscribed to OnCursor, we wouldn't get updated and it would be
		// problematic if t.follow is true. This way we always get updated, if the cursor leaves the forPanel, we
		// unsubscribe from the Manager.win
		Manager().win.SubscribeID(OnCursor, tooltipID, func(s string, i interface{}) {

			if t.subscribed == false {
				return
			}

			if t.status == tooltipActive && t.follow == false {
				return
			}

			// Get info from event
			t.lastevent = i.(*window.CursorEvent)

			if t.status == tooltipActive && t.follow == true {
				t.setLocation()
				return
			}

			if t.status == tooltipPending {
				return
			}

			if t.delay <= 0 {
				t.tooltipDraw(forPanel) // draw now
				return
			}

			// Draw the tooltip after given delay
			t.status = tooltipPending
			go func() {
				select {
				case <-time.After(t.delay):
					t.tooltipDraw(forPanel)
				case <-t.cancel:
					t.status = tooltipOff
				}
			}()
		})

		t.subscribed = true
		Manager().win.Dispatch(OnCursor, i) // Dispatch once so it takes effect
	})

	// When the cursor leaves the forPanel, we cancel or remove the tooltip based on the state
	forPanel.SubscribeID(OnCursorLeave, tooltipID, func(s string, i interface{}) {

		// remove our subscription to Manager.win
		t.subscribed = false
		Manager().win.UnsubscribeAllID(tooltipID)

		switch t.status {
		case tooltipActive:
			t.tooltipClose(forPanel)

		case tooltipPending:
			t.cancel <- true

		case tooltipOff:
		default:
		}
	})
}

// Get the highest Parent in the chain
func (t *Tooltip) getMaster(target IPanel) *core.Node {

	// Recursive Parent() lookup
	var master core.INode = target
	i := 0
	for {
		i++
		if v, ok := master.Parent().(core.INode); ok {

			master = v
			continue
		}
		break
	}
	return master.GetNode()
}

// Draw the tooltip to the screen
func (t *Tooltip) tooltipDraw(target IPanel) {

	t.status = tooltipActive

	// Find and set the position for the tooltip
	t.setLocation()

	// Add the the top most Node
	t.getMaster(target).Add(t.panel)
}

// Set the correct location of the tooltip
func (t *Tooltip) setLocation() {

	if t.align == tooltipFixed {
		return // the t.panel has his position set by t.SetPositionFixed()
	}

	if t.align == tooltipCustom {
		t.panel.SetPosition(t.custom(t.lastevent))
		return
	}

	// Default

	maxw, maxh := Manager().win.GetSize()
	t.panel.SetPosition(t.calculateLocation(t.lastevent.Xpos, t.lastevent.Ypos, float32(maxw), float32(maxh)))
}

// Calculates the correct location for the tooltip
func (t *Tooltip) calculateLocation(mousex, mousey, maxw, maxh float32) (float32, float32) {

	// Get the dimensions of the tooltip
	tooltipHeight := t.panel.Height()
	tooltipWidth := t.panel.Width()

	// Check bottom right
	if (mousey+t.offsetY+tooltipHeight <= maxh) && (mousex+t.offsetX+tooltipWidth <= maxw) {
		return mousex + t.offsetX, mousey + t.offsetY
	}

	// check bottom left
	if mousey+t.offsetY+tooltipHeight <= maxh {
		return (mousex - tooltipWidth) + (t.offsetX * -1), mousey + t.offsetY
	}

	// check top right
	if mousex+t.offsetX+tooltipWidth <= maxw {
		return mousex + t.offsetX, (mousey - tooltipHeight) + (t.offsetY * -1)
	}

	// top left
	return (mousex - tooltipWidth) + (t.offsetX * -1), (mousey - tooltipHeight) + (t.offsetY * -1)
}

// Closes an active tooltip
func (t *Tooltip) tooltipClose(target IPanel) {

	// Remove the panel from the top most node
	if t.getMaster(target).Remove(t.panel) {
		t.status = tooltipOff
	}
}
