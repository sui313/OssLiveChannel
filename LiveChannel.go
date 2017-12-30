package OssLiveChannel

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

//enabled和disabled
func PutLiveChannelStatus(buket, host, channelname, newStatus string) error {
	///ChannelName/PlaylistName?vod&endTime=EndTime&startTime=StartTime
	addr := fmt.Sprintf("http://%s/%s?live&status=%s", host, channelname, newStatus)
	req, err := http.NewRequest(http.MethodPut, addr, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))

	r := &request{}
	r.bucket = buket
	r.path = fmt.Sprintf("/%s", channelname)
	r.headers = req.Header

	para := make(url.Values)
	para.Add("status", newStatus)
	para.Add("live", "")
	r.params = para

	r.method = http.MethodPut
	req.Proto = "HTTP/1.0"

	signRequest(r)
	req.Header = r.headers
	hc := &http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	fmt.Printf("resp:%+v\n", resp)
	buf, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(buf))
	if resp.StatusCode == 200 {
		return nil
	} else {
		return fmt.Errorf("unkown error")
	}
	return nil
}

//创建推流
func PutLiveChannel(data LiveChannelConfiguration, buket, host, channelname string) {
	hc := &http.Client{}

	xmlp := xml.Header
	b, _ := xml.MarshalIndent(data, " ", "  ")
	bb := xmlp + string(b)
	body := bytes.NewReader([]byte(bb))
	addr := fmt.Sprintf("http://%s/%s?live", host, channelname)
	req, err := http.NewRequest(http.MethodPut, addr, body)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))
	//req.Header.Add("Authorization", "OSS "+AccessKeyId+":")
	req.Header.Add("Content-Type", "application/xml")
	req.Header.Add("Content-Length", fmt.Sprintf("%d", len([]byte(bb))))

	h := md5.New()
	io.WriteString(h, bb)
	str_md5 := hex.EncodeToString(h.Sum(nil))
	req.Header.Add("Content-Md5", str_md5)
	//fmt.Println(bb)
	r := &request{}
	r.bucket = buket
	r.path = fmt.Sprintf("/%s?live", channelname)
	r.headers = req.Header
	r.method = "PUT"
	req.Proto = "HTTP/1.0"
	signRequest(r)
	req.Header = r.headers
	resp, err := hc.Do(req)
	if err != nil {
		panic(err)
	}
	//buf := bytes.Buffer
	fmt.Printf("%+v\n", resp)
	fmt.Println("----------")
	buf, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(buf))

	fmt.Println("----------")
	fmt.Println(len(buf))
	fmt.Println("----------")
}

//创建播放列表
func PostVodPlaylist(buket, host, channelname, playlistName, startTime, endTime string) error {
	///ChannelName/PlaylistName?vod&endTime=EndTime&startTime=StartTime
	addr := fmt.Sprintf("http://%s/%s/%s?vod&endTime=%s&startTime=%s",
		host, channelname, playlistName, endTime, startTime)
	req, err := http.NewRequest(http.MethodPost, addr, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))

	para := make(url.Values)
	para.Add("endTime", endTime)
	para.Add("startTime", startTime)
	para.Add("vod", "")

	//req.Header.Add("Content-Length", fmt.Sprintf("%d", len([]byte(bb))))

	r := &request{}
	r.bucket = "app8hyf"
	r.path = fmt.Sprintf("/%s/%s", channelname, playlistName)
	r.headers = req.Header

	r.params = para

	r.method = http.MethodPost
	req.Proto = "HTTP/1.0"

	signRequest(r)
	req.Header = r.headers
	hc := &http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		return nil
	} else {
		fmt.Printf("%+v\n", resp)
		buf, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(buf))
		return fmt.Errorf("unkown error")
	}
	return nil
}

func GetLiveChannelInfo(buket, host, channelname string) {
	GetLiveChannelStat(false, buket, host, channelname)
}

//获取状态
func GetLiveChannelStat(Stat bool, buket, host, channelname string) {
	hc := &http.Client{}

	reqUrl := fmt.Sprintf("http://%s/%s?live", host, channelname)
	if Stat {
		reqUrl = reqUrl + `&comp=stat`
	}

	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))
	//req.Header.Add("Content-Type", "application/xml")
	//req.Header.Add("Content-Length", fmt.Sprintf("%d", len([]byte(bb))))
	para := make(url.Values)
	if Stat {
		para.Add("comp", "stat")
	}
	para.Add("live", "")

	r := &request{}
	r.bucket = buket
	r.path = fmt.Sprintf("/%s", channelname)
	r.params = para
	r.headers = req.Header
	r.method = "GET"
	req.Proto = "HTTP/1.0"
	signRequest(r)
	req.Header = r.headers
	fmt.Println("----------")
	fmt.Println(req.Header.Get("Authorization"))
	fmt.Println(req.Host)
	fmt.Println("----------")
	resp, err := hc.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
	//buf := bytes.Buffer

	if resp.StatusCode == 200 {
		fmt.Println("ok 200")
		//return nil
	} else {
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println("----------")
		fmt.Println(len(buf))
		fmt.Println("----------")
		fmt.Println(string(buf))
		fmt.Println("----------")
		var RespErr LiveRespError
		err = xml.Unmarshal(buf, &RespErr)
		if err != nil {
			panic(err)
		}
		fmt.Println(RespErr.Code)
		fmt.Println(RespErr.ChannelId)
		fmt.Println(RespErr.HostId)
		fmt.Println(RespErr.Message)
		fmt.Println(RespErr.RequestId)
		//return fmt.Errorf("unkown error")
	}

}

func DeleteLiveChannel(buket, host, channelname string) error {
	///ChannelName/PlaylistName?vod&endTime=EndTime&startTime=StartTime
	addr := fmt.Sprintf("http://%s/%s?live", host, channelname)
	req, err := http.NewRequest(http.MethodDelete, addr, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))

	para := make(url.Values)
	para.Add("live", "")

	r := &request{}
	r.bucket = buket
	r.path = fmt.Sprintf("/%s", channelname)
	r.headers = req.Header

	r.params = para

	r.method = http.MethodDelete
	req.Proto = "HTTP/1.0"

	signRequest(r)
	req.Header = r.headers
	hc := &http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}

	fmt.Printf("resp:%+v\n", resp)
	buf, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(buf))

	if resp.StatusCode == 200 {
		return nil
	} else {
		return fmt.Errorf("unkown error")
	}
	return nil
}

func GetLiveChannelHistory(buket, host, channelname string) error {
	//GET /ChannelName?live&comp=history
	addr := fmt.Sprintf("http://%s/%s?live&comp=history", host, channelname)
	req, err := http.NewRequest(http.MethodDelete, addr, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))

	para := make(url.Values)
	para.Add("live", "")
	para.Add("comp", "history")

	r := &request{}
	r.bucket = buket
	r.path = fmt.Sprintf("/%s", channelname)
	r.headers = req.Header

	r.params = para

	r.method = http.MethodGet
	req.Proto = "HTTP/1.0"

	signRequest(r)
	req.Header = r.headers
	hc := &http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}

	fmt.Printf("resp:%+v\n", resp)
	buf, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(buf))

	if resp.StatusCode == 200 {
		var history LiveChannelHistory
		err = xml.Unmarshal(buf, &history)
		if err != nil {
			panic(err)
		}
		return nil
	} else {
		var RespErr LiveRespError
		err = xml.Unmarshal(buf, &RespErr)
		if err != nil {
			panic(err)
		}
		return fmt.Errorf("unkown error")
	}
	return nil
}

func ListLiveChannel(buket, host, marker, prefix string, max_keys int) error {

	para := make(url.Values)
	para2 := make(url.Values)
	para.Add("live", "")
	if marker != "" {
		para.Add("marker", marker)
		para2.Add("marker", marker)
	}
	if prefix != "" {
		para.Add("prefix", prefix)
		para2.Add("prefix", prefix)
	}
	if max_keys > 0 {
		para.Add("max-keys", fmt.Sprintf("%d", max_keys))
		para2.Add("max-keys", fmt.Sprintf("%d", max_keys))
	}
	fmt.Println("===================", para2.Encode())
	addr := fmt.Sprintf("http://%s/?live&%s", host, para2.Encode())
	req, err := http.NewRequest(http.MethodDelete, addr, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))
	r := &request{}
	r.bucket = buket
	r.path = fmt.Sprintf("/")
	r.headers = req.Header

	r.params = para

	r.method = http.MethodGet
	req.Proto = "HTTP/1.0"

	signRequest(r)
	req.Header = r.headers
	hc := &http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	fmt.Println("===================")
	//fmt.Printf("resp:%+v\n", resp)
	buf, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(buf))
	fmt.Println("===================")
	fmt.Println(resp)
	if resp.StatusCode == 200 {
		var Result ListLiveChannelResult
		err = xml.Unmarshal(buf, &Result)
		if err != nil {
			panic(err)
		}
		return nil
	} else {
		var RespErr LiveRespError
		err = xml.Unmarshal(buf, &RespErr)
		if err != nil {
			panic(err)
		}
		return fmt.Errorf("unkown error")
	}
	return nil
}
