{{/*
	Allows you to enlarge virtually any emoji.

	>>> See <https://yagpdb-cc.github.io/utilities/big-emoji-v2> for more information.

	ORIGINAL Author: dvoraknt <https://github.com/dvoraknt>
	LINK TO ORIGINAL CODE: https://github.com/dvoraknt/deebos-yag-CCs/blob/main/on-official-repo/bigEmoji-v2.go.tmpl

	EDITED BY Mei-Jabami: https://github.com/literallymei

	Mei's (me) Note: The original code would send the enlarged emoji PNG as a part of an embed in chat.
	This version of the code would send the emoji itself and that's it. Just what your average user wants.

	Twemoji images are licensed under CC-BY 4.0 and are provided by the official Twemoji project at https://github.com/twitter/twemoji
*/}}

{{deleteTrigger 0}}
{{if .CmdArgs}}
	{{$waitMsg := sendMessageRetID nil (print "<a:loading:844230891781226496> **Loading** <a:loading:844230891781226496>")}}
	{{$embed := sdict}}{{$subArg := index .CmdArgs 0}}{{$emoji := ""}}{{$defEmoji := ""}}{{$fields := cslice}}{{$chan := 0}}{{$msg := 0}}{{$error := false}}{{$ogMsg := ""}}

	{{if reFind `^(\d{17,20})` $subArg}}
		{{if getMessage nil $subArg}}
			{{$msg = $subArg}}{{$chan = .nil}}
			{{$ogMsg = (getMessage nil $subArg).Content}}
			{{$emoji = reFindAllSubmatches `<(a)?:[\w~]+:(\d+)>` $ogMsg}}
			{{$defEmoji = reFindAllSubmatches `([\x{1f1e6}-\x{1f1ff}]{2}|\p{So}\x{fe0f}?[\x{1f3fb}-\x{1f3ff}]?(\x{200D}\p{So}\x{fe0f}?[\x{1f3fb}-\x{1f3ff}]?)*|[#\d*]\x{FE0F}?\x{20E3})` $ogMsg}}
		{{else}}
			{{editMessage nil $waitMsg (print "**Message not found:** If the message isn't in <#" .Channel.ID "> you'll need to use the full message link. Make sure that the message has not been deleted.")}}
			{{deleteMessage nil $waitMsg 10}}
			{{$error = true}}
		{{end}}

	{{else if $linkArg := reFind `(?:[^<]|\A)https?:\/\/(?:www\.)?(?:ptb\.|canary\.)?discord(?:app)?\.com\/channels\/(\d{17,19})\/(\d{17,19})\/(\d{17,19})(?:[^\d]|\z)|<https?:\/\/(?:www\.)?(?:ptb\.|canary\.)?discord(?:app)?\.com\/channels\/(\d{17,19})\/(\d{17,19})\/(\d{17,19})(?:[^>\d]|\z)` $subArg}}
		{{$linkVar := (reFindAll `\d+` $linkArg)}}
		{{$chan = index $linkVar 1}}{{$msg = index $linkVar 2}}
		{{if getMessage $chan $msg}}
			{{$ogMsg = (getMessage $chan $msg).Content}}
			{{$emoji = reFindAllSubmatches `<(a)?:[\w~]+:(\d+)>` $ogMsg}}
			{{$defEmoji = reFindAllSubmatches `([\x{1f1e6}-\x{1f1ff}]{2}|\p{So}\x{fe0f}?[\x{1f3fb}-\x{1f3ff}]?(\x{200D}\p{So}\x{fe0f}?[\x{1f3fb}-\x{1f3ff}]?)*|[#\d*]\x{FE0F}?\x{20E3})` $ogMsg}}
		{{else}}
			{{editMessage nil $waitMsg (print "**Message not found:** Make sure YAGPDB has permission to read messages in <#" $chan "> and that the message has not been deleted.")}}
			{{deleteMessage nil $waitMsg 10}}
			{{$error = true}}
		{{end}}
		
	{{else if eq $subArg "help"}}
		{{$helpEmbed := sdict
		
		"title" (joinStr "" "Big Emoji Help") 
		"description" "This command will allow you to view a single emoji as a larger image or generate image/gif links for up to 25 individual emojis.\n\nYou can use it to view your own emojis, emojis in other messages, used as reactions, or even in a different channel. Use any message ID or message link to extract the emojis and enlarge them.\n\nTo capture reactions use the optional `-re` flag after the message ID or link." 
		"color" 4645612 
		"fields" (cslice 
			(sdict "name" "Syntax" "value" "```elm\n-bigemoji <Emoji> (minimum 1, maximum 25)\n-bigemoji <MessageID> (use when original message exists within the same channel)\n-bigemoji <MessageLink> (use to view reactions from anywhere that YAGPDB has read access)```" "inline" false)
			(sdict "name" "Reaction Flag Usage" "value" "```elm\n-bigemoji <MessageID> -re\n-bigemoji <MessageLink> -re```" "inline" false)
			(sdict "name" "Available Triggers" "value" "`-bigemoji` `-bigemote` `-big-emoji` `-big-emote` `-be`")
			(sdict "name" "Information" "value" "Due to a few limitations some default emojis won't generate proper links or may be incorrect altogether. Most will work but some will not, sorry for any inconvenience.\n\nCustom emojis are displayed in their actual image size, default emojis are displayed in 72x72 as that is largest size the Twemoji CDN provides in PNG format." "inline" false)
		)
		"footer" (sdict "text" .Message.Author.String "icon_url" (.User.AvatarURL "256"))
		"timestamp" currentTime
		}}
		{{editMessage nil $waitMsg (complexMessageEdit "content" "" "embed" (cembed $helpEmbed))}}
		
	{{else if eq $subArg "-re"}}
		{{editMessage nil $waitMsg (print "**Invalid Syntax:** The `-re` flag must be placed after the message ID/link.")}}
		{{deleteMessage nil $waitMsg 10}}
		{{$error = true}}

	{{else}}
		{{$emoji = reFindAllSubmatches `<(a)?:[\w~]+:(\d+)>` .Message.Content}}
		{{$defEmoji = reFindAllSubmatches `([\x{1f1e6}-\x{1f1ff}]{2}|\p{So}\x{fe0f}?[\x{1f3fb}-\x{1f3ff}]?(\x{200D}\p{So}\x{fe0f}?[\x{1f3fb}-\x{1f3ff}]?)*|[#\d*]\x{FE0F}?\x{20E3})` .Message.Content}}
	{{end}}
	
	{{if and (reFind `(?i)-re` (print .CmdArgs)) $ogMsg}}
	{{$emoji = cslice.AppendSlice $emoji}}{{$defEmoji = cslice.AppendSlice $defEmoji}}
		{{range (getMessage $chan $msg).Reactions}}
			{{if and (ne (toInt .Emoji.ID) 0) (not .Emoji.Animated)}}
				{{- $emoji = $emoji.AppendSlice (cslice (cslice (printf "<:%s:%d>" .Emoji.Name .Emoji.ID) "" .Emoji.ID)) -}}
			{{else if (ne (toInt .Emoji.ID) 0) .Emoji.Animated}}
				{{- $emoji = $emoji.AppendSlice (cslice (cslice (printf "<:%s:%d>" .Emoji.Name .Emoji.ID) "a" .Emoji.ID)) -}}
			{{else if eq (toInt .Emoji.ID) 0}}
				{{- $defEmoji = $defEmoji.AppendSlice (cslice (cslice .Emoji.Name)) -}}
			{{end}}
		{{end}}
	{{end}}
	
	{{if $emoji}}
	{{if eq (add (len $emoji) (len $defEmoji)) 1}}
		{{with $emoji}}
			{{$animated := index . 0 1}}
			{{$id := str (index . 0 2)}}
			{{$ext := ".png"}}{{$typeExt := "Image"}}
			{{if $animated}}{{$ext = ".gif"}} {{$typeExt = "Gif"}}{{end}}
			{{(print "https://cdn.discordapp.com/emojis/" $id $ext)}}
			{{deleteMessage nil $waitMsg 0}}
		{{end}}
	{{else}}
		{{range $emoji}}
			{{$animated := index . 1}}
			{{$id := index . 2}}
			{{$ext := ".png"}}{{$typeExt := "Image"}}
			{{if $animated}}{{$ext = ".gif"}}{{$typeExt = "Gif"}}{{end}}
			{{(print "[" $typeExt " Link](https://cdn.discordapp.com/emojis/" $id $ext ")")}}
		{{end}}
		{{$embed.Set "title" "This message contains more than one emoji."}}
		{{$embed.Set "fields" $fields}}
	{{end}}
	{{end}}
	
	{{if $defEmoji}}
		{{$emoji_U := ""}}{{$url := "https://twemoji.maxcdn.com/v/latest/72x72/"}}

		{{if eq (add (len $emoji) (len $defEmoji)) 1}}
			{{- range toRune (index $defEmoji 0 0)}}
				{{- $emoji_U = joinStr "-" $emoji_U (printf "%04x" .)}}
			{{- end -}}
			{{$url = joinStr "" $url $emoji_U ".png"}}
			{{print (joinStr "" $url)}}
			{{deleteMessage nil $waitMsg 0}}
		{{else}}
			{{range $defEmoji}}
				{{- range toRune (index . 0)}}
					{{- $emoji_U = joinStr "-" $emoji_U (printf "%04x" . )}}
				{{end}}
				{{- $fields = $fields.Append (sdict "name" (print "`" (index . 0) "`") "value" (print "[Image Link](" (joinStr "" $url $emoji_U ".png") ")") "inline" true) -}}
				{{$emoji_U = ""}}
			{{end}}
			{{$embed.Set "title" "This message contains more than one emoji."}}
			{{$embed.Set "fields" $fields}}
		{{end}}
	{{end}}
		
	{{if and (not $emoji) (not $defEmoji) (not $error) (not (eq $subArg "help"))}}
		{{editMessage nil $waitMsg (print "This message does not contain any emojis or you have given an incorrect message ID.")}}
		{{deleteMessage nil $waitMsg 10}}
	{{else if and (gt (len $fields) 25) (not $error) (not (eq $subArg "help"))}}
		{{editMessage nil $waitMsg (print "There are more than 25 emojis in the message, please try again with fewer emojis or with a different source message.")}}
		{{deleteMessage nil $waitMsg 10}}
	{{else if $embed}}
		{{$embed.Set "color" 0x39ff14}}
		{{$embed.Set "title" (joinStr "" "Big Emoji")}}
		{{if $defEmoji}}
			{{$embed.Set "footer" (sdict "text" (print "Default emoji not correct? Check help.\n" .Message.Author.String "  |  -bigemoji help") "icon_url" (.User.AvatarURL "256"))}}
		{{else}}
			{{$embed.Set "footer" (sdict "text" (print .Message.Author.String "  |  -bigemoji help") "icon_url" (.User.AvatarURL "256"))}}
		{{end}}
		{{editMessage nil $waitMsg (complexMessageEdit "content" "" "embed" (cembed $embed))}}
	{{end}}

{{else}}
	{{$failMsg := sendMessageRetID nil (print "No arguments provided! Use `-bigemoji help` for information on how to use this command.")}}
	{{deleteMessage nil $failMsg 10}}
{{end}}
