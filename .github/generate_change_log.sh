#!/usr/bin/env bash
checksum() {
    echo $(sha256sum $@ | awk '{print $1}')
}
change_log_file="./CHANGELOG.md"
version="## $@"
version_prefix="## v"
start=0
CHANGE_LOG=""
while read line; do
    if [[ $line == *"$version"* ]]; then
        start=1
        continue
    fi
    if [[ $line == *"$version_prefix"* ]] && [ $start == 1 ]; then
        break;
    fi
    if [ $start == 1 ] && [[ $line != "" ]]; then
        CHANGE_LOG+="$line\n"
    fi
done < ${change_log_file}
LINUX_BIN_SUM="$(checksum ./linux/kukd)"
MAC_BIN_SUM="$(checksum ./macos/kuskd)"
WINDOWS_BIN_SUM="$(checksum ./windows/kuskd)"
OUTPUT=$(cat <<-END
## Changelog\n
${CHANGE_LOG}\n
## Assets\n
|    Assets    | Sha256 Checksum  |\n
| :-----------: |------------|\n
| kuskd_linux | ${LINUX_BIN_SUM} |\n
| kuskd_mac  | ${MAC_BIN_SUM} |\n
| kuskd_windows  | ${WINDOWS_BIN_SUM} |\n
END
)

echo -e ${OUTPUT}