  TODO
├╴󰟓 internal/http/handlers/channel-handlers.go 5
│ ├╴ TODO make a better struct for all [70, 5]
│ ├╴ TODO no need for this [155, 6]
│ ├╴ TODO this function is a mess [335, 4]
│ ├╴ TODO add logic that checks if the user is an owner of this channel [348, 5]
│ └╴ TODO send this message to the user [365, 26]
├╴󰌜 assets/css/buttons.css 1
│ └╴ TODO embed SVGs*/ [202, 3]
├╴󰌜 assets/css/filters.css 1
│ └╴ TODO add css for days if they show for previosu/next month in left/right calendars*/ [329, 40]
├╴󰌜 assets/css/main.css 8
│ ├╴ TODO ---------- convert all px to rem ---------- _/ [15, 5]
│ ├╴ TODO add a container query to size this correctly _/ [266, 6]
│ ├╴ TODO use @starting-style to achieve this _/ [988, 4]
│ ├╴ TODO prevent this happening on form popop _/ [1000, 5]
│ ├╴ TODO get code under working _/ [1198, 6]
│ ├╴ TODO get code under working _/ [1214, 6]
│ ├╴ TODO does this need a var(--cursor-default)? [1460, 4]
│ └╴ TODO make this transition work _!_/ [2069, 10]
├╴󰜡 assets/images/Facebook_logo.svg 1
│ └╴ TODO check https works--> [30, 22]
├╴󰌞 assets/js/authentication.js 2
│ ├╴ TODO overhaul the naming of these buttons [24, 4]
│ └╴ TODO get these working [122, 4]
├╴󰌞 assets/js/calendar.js 1
│ └╴ TODO need to add day text for the days of the previous month [316, 6]
├╴󰌞 assets/js/fetch_and_navigate.js 3
│ ├╴ TODO need to modify home-page template to populate by injection [55, 10]
│ ├╴ TODO when injected, the content can be cleared [56, 10]
│ └╴ TODO with injected home-page, this can be removed [61, 10]
├╴󰌞 assets/js/filters.js 4
│ ├╴ TODO commenting out type as not needed for base audit [109, 8]
│ ├╴ TODO commenting out type as not needed for base audit [120, 12]
│ ├╴ TODO commenting out type as not needed for base audit [195, 12]
│ └╴ TODO commenting out type as not needed for base audit [209, 12]
├╴󰌞 assets/js/popups.js 1
│ └╴ TODO refactor the open and close modals [62, 6]
├╴󰌞 assets/js/share.js 8
│ ├╴ TODO Add logic that positions the modal above the button if there's not enough space under [24, 4]
│ ├╴ TODO check api's of sites and fix title/message [137, 8]
│ ├╴ TODO check https works [154, 10]
│ ├╴ TODO check https works [157, 10]
│ ├╴ TODO check https works [160, 10]
│ ├╴ TODO check https works [163, 10]
│ ├╴ TODO check https works [166, 10]
│ └╴ TODO check https works [169, 10]
├╴󰬁 assets/templates/channel-edit-mods.tmpl 1
│ └╴ TODO range through moderators _/}} [25, 14]
├╴󰬁 assets/templates/control-buttons.tmpl 2
│ ├╴ TODO - simply don't add this template to right-panel, omit this if clause _/}} [30, 6]
│ └╴ TODO when user is an administrator, button prepared; needs admin check*/}} [109, 6]
├╴󰬁 assets/templates/filters-row.tmpl 4
│ ├╴ TODO add hide-feed for all but home-page activity */}} [12, 7]
│ ├╴ TODO commented out until further release as not needed for base audit _/}} [46, 13]
│ ├╴ TODO commented out until further release as not needed for base audit _/}} [107, 13]
│ └╴ TODO calendar setup for the additional features _/}} [214, 9]
├╴󰬁 assets/templates/list-item.tmpl 1
│ └╴ TODO add links to sidebar elements_/}} [12, 7]
├╴󰍔 audit/base_audit.md 6
│ ├╴ TODO** we need to sort out the pages [213, 3]
│ ├╴ TODO** we need to add this functionality (issue \#65) [219, 3]
│ ├╴ TODO** we need to handle the redirects correctly (issue \#66) [243, 3]
│ ├╴ TODO** we need to explicitly handle this (issue \#67) [247, 3]
│ ├╴ TODO** we need to explicitly handle this (issue \#67) [251, 3]
│ └╴ TODO** create the build file (and check the build) [267, 3]
├╴󰟓 cmd/server/main.go 1
│ └╴ TODO figure this out [55, 5]
├╴󰟓 internal/http/handlers/user-handlers.go 1
│ └╴ TODO does this check need to be here? [226, 5]
├╴󰟓 internal/models/rule-models.go 1
│ └╴ TODO figure out the use of strings/int64s for this [18, 4]
├╴󰟓 internal/sqlite/channels-sql.go 4
│ ├╴ TODO (realtime) get this data from websockets [55, 6]
│ ├╴ TODO (realtime) get this data from websockets [105, 6]
│ ├╴ TODO (realtime) get this data from websockets [140, 6]
│ └╴ TODO (realtime) get this data freom websockets [210, 6]
├╴󰟓 internal/sqlite/comments-sql.go 1
│ └╴ TODO add Updated field, which should be populated on update [112, 5]
├╴󰟓 internal/sqlite/posts_test-sql.go 3
│ ├╴ TODO: Add test cases. [21, 6]
│ ├╴ TODO: Add test cases. [54, 6]
│ └╴ TODO: Add test cases. [93, 6]
├╴󰟓 internal/sqlite/reactions-sql.go 1
│ └╴ TODO refactor so that query inserts ID/NULL to PostID AND CommentID [102, 5]
├╴󰟓 internal/sqlite/reactions_test-sql.go 1
│ └╴ TODO: Add test cases. [27, 6]
└╴󰟓 internal/sqlite/users-sql.go 2
├╴ TODO unify these functions to accept parameters [156, 4]
└╴ TODO accept an interface for any given value [258, 4]
  FIX
├╴󰌜 assets/css/buttons.css 1
│ └╴ FIXME span css conflicts with .btn-icotext & span &::before _/ [334, 4]
├╴󰬁 assets/templates/channel-join-popover.tmpl 1
│ └╴ FIXME ---------- fix these buttons ----------_/}} [31, 5]
├╴󰬁 assets/templates/channel-page-banner.tmpl 1
│ └╴ FIXME base audit - commented out until all buttons work _/}} [29, 12]
├╴󰌝 assets/templates/index.html 2
│ ├╴ FIXME ---------- fix resize issue ---------- _/}} [341, 7]
│ └╴ FIXME ---------- fix these buttons ---------- _/}} [363, 7]
├╴󰬁 assets/templates/post-card.tmpl 1
│ └╴ FIXME base audit - removed links and data-dests for user-page being commented out_/}} [27, 21]
├╴󰟓 internal/sqlite/channels-sql.go 1
│ └╴ FIXME: This is a temporary fix to set the channel as joined:we need to come up with a more robust solution [53, 6]
└╴󰟓 internal/sqlite/users-sql.go 2
├╴ FIXME this prepare statement is unnecessary as it is not used in a loop [31, 5]
└╴ FIXME this error [212, 7]
