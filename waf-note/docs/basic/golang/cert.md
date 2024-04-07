## 背景
Golang验证证书公私钥的有效性


## 相关文章：
https://iwiki.woa.com/pages/viewpage.action?pageId=1244985579
https://developer.aliyun.com/article/284911


## encoding/pem
https://studygolang.com/static/pkgdoc/pkg/encoding_pem.htm
Decode函数会从输入里查找到下一个PEM格式的块（证书、私钥等）。它返回解码得到的Block和剩余未解码的数据。如果未发现PEM数据，返回(nil, data)。


### 解码公钥
证书公钥存在证书链的情况：

``` 
-----BEGIN CERTIFICATE-----
......
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
......
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
......
-----END CERTIFICATE-----

```

因此decode的时候，需要循环解析证书链：
``` 
certRawList := make([][]byte, 0)
for {
	pemBlock, rest := pem.Decode(certRaw)
	if pemBlock != nil {
		certRawList = append(certRawList, pemBlock.Bytes)
	}
}
```

### 验证公钥

``` 
读取系统的根证书

var rootCertPool *x509.CertPool

//InitRootCertPool 初始化根证书目录
func InitRootCertPool(rootCertPath string) error {
    var err error
    rootCertPool, err = x509.SystemCertPool()
    if err != nil {
        return err
    }

    fis, err := ioutil.ReadDir(rootCertPath)
    if err != nil {
        return err
    }
    for _, fi := range fis {
        data, err := ioutil.ReadFile(rootCertPath + "/" + fi.Name())
        if err != nil {
            return err
        }
        if !rootCertPool.AppendCertsFromPEM(data) {
            log.Logger.Errorf("append cert %s failed", fi.Name())
        }
        log.Logger.Infof("append cert %s success", fi.Name())
    }
    return nil
}
```


``` 
// 基于根证书 验证用户证书
	
	//区分用户证书和中间证书
	Intermediates := x509.NewCertPool() //中间证书
    for _, certRawData := range certRawList {
        certInfo, err := x509.ParseCertificate(certRawData)
        if err != nil {
          	//解析错误
        }

        if t.ParsedCert != nil {
            Intermediates.AddCert(certInfo)//后面的证书都加入中间证书
        } else {
            t.ParsedCert = certInfo //第一个证书为当前要验证的这个数
        }
    }
	
	
	// 基于中间证书和根证书，开始验证用户证书的有效性
	_, err := t.ParsedCert.Verify(x509.VerifyOptions{
        Intermediates: Intermediates,
        Roots:         rootCertPool,
    })
	if err != nil {
        // 证书无效，CA根证书验证失败
    }
	
	// 验证过期时间
	if t.ParsedCert.NotAfter.Sub(now).Hours() <= 15*24 { // 15天*24小时
        errmsg := fmt.Sprintf("certificate will expire at: %s", t.ParsedCert.NotAfter.String())
    }
	
	// 验证域名信息
    if err := t.ParsedCert.VerifyHostname(t.Domain); err != nil {
        logger.Errorf("VerifyHostname failed: %s", err.Error())
    }
```





### 解码私钥
私钥不存在在链的情况，直接解码即可
``` 
keyDERBlock, _ := pem.Decode(keyRaw)
if keyDERBlock == nil {
	// 解码错误
}

if keyDERBlock.Type == "RSA PRIVATE KEY" {
        //RSA PKCS1
        key, errParsePK = x509.ParsePKCS1PrivateKey(keyDERBlock.Bytes)
    } else if keyDERBlock.Type == "PRIVATE KEY" {
        //pkcs8格式的私钥解析
        key, errParsePK = x509.ParsePKCS8PrivateKey(keyDERBlock.Bytes)
    }

if errParsePK != nil {
	return nil // 解析错误
} else {
	cert.PrivateKey = key
}
```



