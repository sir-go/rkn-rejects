# get_rkn
  
[CryptoPRO](https://www.cryptopro.ru/products/csp/downloads) for sign requests is required.

SSH tuning `.ssh/config`
```
KexAlgorithms +diffie-hellman-group-exchange-sha1,diffie-hellman-group14-sha1
```
## What it does
 - check versions of the service, dump and documentation
 - generate XML request, sign it and send to the service
 - get the request UUID and wait the result
 - download the dump, unpack and smart parse it
 - include own denied lists
 - fill the redis DB
 - generate and run script for configuring the firewall

## Config
```toml
log_level = "debug"

[actual_versions]                       # actual RKN service versions
    service = "3.2"                     # SOAP service
    dump    = "2.4"                     # register dump
    doc     = "4.12"                    # version of the documentation

[web]                                   # SOAP service params
    soap_url    = "https://.../?wsdl"   # WSDL URL
    doc_url     = "https://.../.pdf"    # documentation URL
    tcp_timeout = "30s"                 # TCP timeout
    attempts    = 5                     # amount of connection attempts

[req]                                   # request's prarmeters
    file = "/tmp/req.xml"               # where to save temp. request file 
    [req.operator]
        name    = 'company-name'        # company name
        inn     = "company-inn"         # company TIN
        ogrn    = "company-ogrn"        # company PSRN
        email   = "company@email.com"   # company registered e-mail

[sign]                                  # digital signature 
    file = "/tmp/req.xml.signed"        # where to save temp. sign file
    script = "/opt/rkn/sign.sh"         # path to the signing script

[res]                                   # getting req. result params
## optional, dump will be 
## saved to file if it set
#    dump_to = "/tmp/dump.zip"          # where to save dump file if needed
    attempts            = 30            # amount of attemps to get the result
    retry_timeout       = "15s"         # timeout for each try
    download_timeout    = "10m"         # duration of waiting

[parse]                                 # dump parsing params
## optional, the parser will read this 
## file if it set instead of memory
#    from_dump = "/tmp/dump.zip"        # downloaded dump instead of getting the new one
    progress_poll_timeout = "1s"        # parsing goroutine polling timeout

    [parse.bogus_ip]                    # some CIDRs in the dump are bogus or too wide
        subnets = [                     # array of private CIDRs 
            "0.0.0.0/8",      
            "10.0.0.0/8",     
            "14.0.0.0/8",     
            "169.254.0.0/16", 
            "127.0.0.0/8",    
            "192.168.0.0/16", 
            "172.16.0.0/12",  
            "192.0.2.0/24",   
            "224.0.0.0/3",    
        ]
        min_mask = 10                   # minimal subnet mask (too wide rejects prevention)

    [parse.redis]                       # redis cache connection params
        host            = "localhost"   # redis host
        port            = 6379          # TCP port
        db              = 0             # num of DB
        chunk_size      = 1000          # inserted data chunk size
        workers         = 16            # amount of inserters
        timeout_conn    = "15s"         # connection timeout
        timeout_read    = "15s"         # read processes timeout

[lists]                                             # custom lists
    black_domains   = "wb_lists/b_domains.yml"      # path to the own denied domains list
    white_domains   = "wb_lists/w_domains.yml"      # path to the own allowed domains list

[fw]                                                # firewall params
    ip_deny_file = "/etc/nftables.d/deny_rkn.nft"   # path to save parsed addresses
    ip_deny_table = "rkn"                           # nftables table name
    ip_deny_set = "deny_rkn"                        # nftables denied addresses set name
```

## Build
```bash
go mod download
go build -o get_rkn rkn-rejects/cmd/get_rkn;
```
## Run
```bash
./get_rkn -c rkn.toml
```
## Flags
`-c <config file path>` - path to `*.toml` config file (default is `rkn.toml`)
