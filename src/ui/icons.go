package ui

import "fyne.io/fyne/v2"

// Star icons as embedded SVG resources (24x24).
var starOutlineSVG = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24">
  <path fill="none" stroke="#888" stroke-width="2" d="M12 3.5l2.9 5.88 6.5.95-4.7 4.58 1.1 6.43L12 18.9 6.2 21.34l1.1-6.43-4.7-4.58 6.5-.95L12 3.5z"/>
  <path fill="none" d="M0 0h24v24H0z"/>
 </svg>`)

var starFilledSVG = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24">
  <path fill="#888" d="M12 2.5l3.24 6.57 7.26 1.06-5.25 5.1 1.24 7.22L12 19.77 5.51 22.45l1.24-7.22L1.5 10.13l7.26-1.06L12 2.5z"/>
  <path fill="none" d="M0 0h24v24H0z"/>
 </svg>`)

var StarOutlineIcon fyne.Resource = fyne.NewStaticResource("star-outline.svg", starOutlineSVG)
var StarFilledIcon fyne.Resource = fyne.NewStaticResource("star-filled.svg", starFilledSVG)

// Blue filled star
var starBlueSVG = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24">
  <path fill="#1976D2" d="M12 2.5l3.24 6.57 7.26 1.06-5.25 5.1 1.24 7.22L12 19.77 5.51 22.45l1.24-7.22L1.5 10.13l7.26-1.06L12 2.5z"/>
  <path fill="none" d="M0 0h24v24H0z"/>
</svg>`)
var StarBlueIcon fyne.Resource = fyne.NewStaticResource("star-blue.svg", starBlueSVG)

// Red cross (X) icon as SVG (24x24) with red strokes
var redCrossSVG = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24">
  <line x1="6" y1="6" x2="18" y2="18" stroke="#E53935" stroke-width="2" stroke-linecap="round" />
  <line x1="18" y1="6" x2="6" y2="18" stroke="#E53935" stroke-width="2" stroke-linecap="round" />
  <path fill="none" d="M0 0h24v24H0z"/>
</svg>`)

var RedCrossIcon fyne.Resource = fyne.NewStaticResource("red-cross.svg", redCrossSVG)

// Arrow up/down icons for spinners (16x16)
var arrowUpSVG = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16">
  <path d="M8 3l5 6H9v4H7V9H3l5-6z" fill="#666"/>
  <path fill="none" d="M0 0h16v16H0z"/>
</svg>`)

var arrowDownSVG = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16">
  <path d="M8 13l-5-6h4V3h2v4h4l-5 6z" fill="#666"/>
  <path fill="none" d="M0 0h16v16H0z"/>
</svg>`)

var ArrowUpIcon fyne.Resource = fyne.NewStaticResource("arrow-up.svg", arrowUpSVG)
var ArrowDownIcon fyne.Resource = fyne.NewStaticResource("arrow-down.svg", arrowDownSVG)
