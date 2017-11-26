$gover = "1.9.2"

from "xena/go-mini:#{$gover}"

run "go#{$gover} download"

copy "./cmd/disrelay/main.go", "/root/go/src/github.com/Xe/disrelay/cmd/disrelay/main.go"
copy "./vendor", "/root/go/src/github.com/Xe/disrelay/vendor"

run "go#{$gover} build -o /usr/local/bin/disrelay github.com/Xe/disrelay/cmd/disrelay"

run "apk del go#{$gover} && rm -rf /root/go /root/sdk"

flatten

cmd "/usr/local/bin/disrelay"

tag "xena/disrelay:0.1"
