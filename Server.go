package main
import  ("github.com/julienschmidt/httprouter"
		"fmt"
		"net/http"
		"strconv"
		"encoding/json"
		"strings"
		"sort")


type KeyValue struct{
	Key int	`json:"key,omitempty"`
	Value string	`json:"value,omitempty"`
} 

type CountElements struct{
	Node string   `json:"node"`
	Count int	  `json:"count"`
}

type Mapping struct{
	Node string `json:"node"`	
	Key int `json:"key"`
}

var s1,s2,s3 [] KeyValue
var index1,index2,index3 int
type ByKey []KeyValue
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

type ForMapping []Mapping
func (a ForMapping) Len() int           { return len(a) }
func (a ForMapping) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ForMapping) Less(i, j int) bool { return a[i].Key < a[j].Key }

func GetAllKeys(rw http.ResponseWriter, request *http.Request,p httprouter.Params){
	var out []KeyValue
	for _,i:=range s1{
		out = append(out,i)
	}
	for _,i:=range s2{
		out = append(out,i)
	}
	for _,i:=range s3{
		out = append(out,i)
	}
	sort.Sort(ByKey(out))
	result,_:= json.Marshal(out)
	fmt.Fprintln(rw,string(result))
}

func GetSpecificPortKeys(rw http.ResponseWriter, request *http.Request,p httprouter.Params){
	port := strings.Split(request.Host,":")
	if(port[1]=="3000"){
		sort.Sort(ByKey(s1))
		result,_:= json.Marshal(s1)
		fmt.Fprintln(rw,string(result))
	}else if(port[1]=="3001"){
		sort.Sort(ByKey(s2))
		result,_:= json.Marshal(s2)
		fmt.Fprintln(rw,string(result))
	}else{
		sort.Sort(ByKey(s3))
		result,_:= json.Marshal(s3)
		fmt.Fprintln(rw,string(result))
	}
}

func PutKeys(rw http.ResponseWriter, request *http.Request,p httprouter.Params){
	port := strings.Split(request.Host,":")
	key,_ := strconv.Atoi(p.ByName("key_id"))
	if(port[1]=="3000"){
		s1 = append(s1,KeyValue{key,p.ByName("value")})
		index1++
	}else if(port[1]=="3001"){
		s2 = append(s2,KeyValue{key,p.ByName("value")})
		index2++
	}else{
		s3 = append(s3,KeyValue{key,p.ByName("value")})
		index3++
	}	
}

func GetOneKey(rw http.ResponseWriter, request *http.Request,p httprouter.Params){	
	out := s1
	ind := index1
	port := strings.Split(request.Host,":")
	if(port[1]=="3001"){
		out = s2 
		ind = index2
	}else if(port[1]=="3002"){
		out = s3
		ind = index3
	}	
	key,_ := strconv.Atoi(p.ByName("key_id"))
	for i:=0 ; i< ind ;i++{
		if(out[i].Key==key){
			result,_:= json.Marshal(out[i])
			fmt.Fprintln(rw,string(result))
		}
	}
}


func GetCount(rw http.ResponseWriter, request *http.Request,p httprouter.Params){	
	var output [3]CountElements
	output[0].Node = "http://localhost:3000"
	output[1].Node = "http://localhost:3001"
	output[2].Node = "http://localhost:3002"
	output[0].Count = index1
	output[1].Count = index2
	output[2].Count = index3
	result,_:= json.Marshal(output)
	fmt.Fprintln(rw,string(result))
}

func GetMapping(rw http.ResponseWriter, request *http.Request,p httprouter.Params){
	var output[] Mapping
	var j,k,l,stop1,stop2,stop3 int
	j=0
	k=0
	l=0 
	stop1=0
	stop2=0
	stop3=0
	for i:=0;i<(index1+index2+index3); i++{
		if(s1!=nil && stop1!=1){
			output = append(output,Mapping{"http://localhost:3000",s1[j].Key})
			if(j < len(s1)-1){
				j++
			}else{
				stop1=1
			}
		}
		if(s2!=nil && stop2!=1){
			output = append(output,Mapping{"http://localhost:3001",s2[k].Key})
			if(k < len(s2)-1){
				k++
			}else{
				stop2 =1
			}
		}
		if(s3!=nil && stop3!=1){
			output = append(output,Mapping{"http://localhost:3002",s3[l].Key})
			if(l < len(s3)-1){
				l++
		    }else{
		    	stop3=1
		    }
		}
	}
	sort.Sort(ForMapping(output))
	result,_:= json.Marshal(output)
	fmt.Fprintln(rw,string(result))
}

func RemoveNode(rw http.ResponseWriter, request *http.Request,p httprouter.Params){
	
	var out []KeyValue
	port := strings.Split(request.Host,":")
	if(port[1]=="3000"){
		for _,i:=range s1{
			out = append(out,i)
		}
		s1=nil
		index1=0
	}else if(port[1]=="3001"){
		for _,i:=range s2{
			out = append(out,i)
		}
		s2=nil
		index2=0
	}else if(port[1]=="3002"){
		for _,i:=range s3{
			out = append(out,i)
		}
		s3=nil
		index3=0
	}
	sort.Sort(ByKey(out))
	result,_:= json.Marshal(out)
	fmt.Fprintln(rw,string(result))
}

func main(){
	index1 = 0
	index2 = 0
	index3 = 0
	mux := httprouter.New()
    mux.GET("/keys",GetSpecificPortKeys)
    mux.GET("/keys/:key_id",GetOneKey)
    mux.PUT("/keys/:key_id/:value",PutKeys)
    mux.GET("/count",GetCount)
    mux.GET("/mapping",GetMapping)
    mux.GET("/remove",RemoveNode)
    mux.GET("/all",GetAllKeys)
    go http.ListenAndServe(":3000",mux)
    go http.ListenAndServe(":3001",mux)
    go http.ListenAndServe(":3002",mux)
    select {}
}