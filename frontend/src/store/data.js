import { defineStore } from "pinia";
import { rsgData } from '../utils/data';
const FilterOperatingSystemType = {
    'All': 'all',
    'Windows': 'windows',
    'Linux': 'linux',
    'Mac': 'mac'
};

const hoaxshell_listener_types = {
    "Windows CMD cURL": "cmd-curl",
    "PowerShell IEX": "ps-iex",
    "PowerShell IEX Constr Lang Mode": "ps-iex-cm",
    "PowerShell Outfile": "ps-outfile",
    "PowerShell Outfile Constr Lang Mode": "ps-outfile-cm",
    "Windows CMD cURL https": "cmd-curl -c /your/cert.pem -k /your/key.pem",
    "PowerShell IEX https": "ps-iex -c /your/cert.pem -k /your/key.pem",
    "PowerShell IEX Constr Lang Mode https": "ps-iex-cm -c /your/cert.pem -k /your/key.pem",
    "PowerShell Outfile https": "ps-outfile -c /your/cert.pem -k /your/key.pem",
    "PowerShell Outfile Constr Lang Mode https": "ps-outfile-cm -c /your/cert.pem -k /your/key.pem"
};
const parsePortOrDefault = function (value, defaultPort = 9001) {
    if (value === null || value === undefined) return defaultPort;

    const number = Number(value);
    const isValidPort = (Number.isSafeInteger(number) && number >= 0 && number <= 65535);
    return isValidPort ? number : defaultPort;
};
const query = new URLSearchParams(location.hash.substring(1));
export const useDataStore = defineStore('Data', () => {
    const rsg = {
        ip: (query.get('ip') || localStorage.getItem('ip') || '10.10.10.10').replace(/[^a-zA-Z0-9.\-]/g, ''),
        port: parsePortOrDefault(query.get('port') || localStorage.getItem('port')),
        payload: query.get('payload') || localStorage.getItem('payload') || 'windows/x64/meterpreter/reverse_tcp',
        payload: query.get('type') || localStorage.getItem('type') || 'cmd-curl',
        shell: query.get('shell') || localStorage.getItem('shell') || rsgData.shells[0],
        listener: query.get('listener') || localStorage.getItem('listener') || rsgData.listenerCommands[0][1],
        encoding: query.get('encoding') || localStorage.getItem('encoding') || 'None',
        getIP: () => rsg.ip,

        getPort: () => parsePortOrDefault(rsg.port),

        getShell: () => rsg.shell,

        getEncoding: () => rsg.encoding,

        getSelectedCommandName: () => {
            return rsg.selectedValues[rsg.commandType];
        },

        getReverseShellCommand: () => {
            const reverseShellData = rsgData.reverseShellCommands.find((item) => item.name === rsg.getSelectedCommandName());
            return reverseShellData.command;
        },

        getPayload: () => {
            if (rsg.commandType === 'MSFVenom') {
                let cmd = rsg.getReverseShellCommand();
                // msfvenom -p windows/x64/meterpreter_reverse_tcp ...
                let regex = /\s+-p\s+(?<payload>[a-zA-Z0-9/_]+)/;
                let match = regex.exec(cmd);
                if (match) {
                    return match.groups.payload;
                }
            }

            return 'windows/x64/meterpreter/reverse_tcp'

        },

        getType: () => {
            if (rsg.commandType === 'HoaxShell') {
                let cmd_name = rsg.getSelectedCommandName();
                return hoaxshell_listener_types[cmd_name];
            }

            return 'cmd-curl'

        },

        generateReverseShellCommand: () => {
            let command

            if (rsg.getSelectedCommandName() === 'PowerShell #3 (Base64)') {
                const encoder = (text) => text;
                const payload = rsg.insertParameters(rsgData.specialCommands['PowerShell payload'], encoder)
                command = "powershell -e " + btoa(toBinary(payload))
                function toBinary(string) {
                    const codeUnits = new Uint16Array(string.length);
                    for (let i = 0; i < codeUnits.length; i++) {
                        codeUnits[i] = string.charCodeAt(i);
                    }
                    const charCodes = new Uint8Array(codeUnits.buffer);
                    let result = '';
                    for (let i = 0; i < charCodes.byteLength; i++) {
                        result += String.fromCharCode(charCodes[i]);
                    }
                    return result;
                }
            } else {
                command = rsg.getReverseShellCommand()
            }

            const encoding = rsg.getEncoding();
            if (encoding === 'Base64') {
                command = rsg.insertParameters(command, (text) => text)
                command = btoa(command)
            } else {
                function encoder(string) {
                    let result = string;
                    switch (encoding) {
                        case 'encodeURLDouble':
                            result = fixedEncodeURIComponent(result);
                        // fall-through
                        case 'encodeURL':
                            result = fixedEncodeURIComponent(result);
                            break;
                    }
                    return result;
                }
                command = rsg.escapeHTML(encoder(command));
                // NOTE: Assumes encoder doesn't produce HTML-escaped characters in parameters
                command = rsg.insertParameters(rsg.highlightParameters(command, encoder), encoder);
            }

            return command;
        },
        fixedEncodeURIComponent: (str) => {
            return encodeURIComponent(str).replace(/[!'()*]/g, function (c) {
                return '%' + c.charCodeAt(0).toString(16).toUpperCase();
            });
        },
        highlightParameters: (text, encoder) => {
            // 确保 text 是一个字符串
            text = String(text);

            const parameters = ['{ip}', '{port}', '{shell}', encodeURI('{ip}'), encodeURI('{port}'), encodeURI('{shell}')];
            // console.log(parameters)

            parameters.forEach((param) => {
                if (encoder) param = encoder(param);
                // console.log(param)
                text = text.replace(new RegExp(param, 'g'), `<span class="highlighted-parameter">${param}</span>`);

            });

            return text;
        },
    }
    function portImcrement() {
        this.rsg.port++
    }
    return { rsg, portImcrement }
})