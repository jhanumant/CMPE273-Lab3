package main  

  
import (  
    "fmt"  
    "hash/crc32"  
    "sort"     
    "net/http"
    "encoding/json" 
    "io/ioutil"
    "strconv"
)  
   
type HashCircle []uint32  

type KeyValue struct{
    Key int `json:"key,omitempty"`
    Value string `json:"value,omitempty"`
}

type Mapping struct{
    NodeAddress string `json:"node,omitempty"`   
    Key int `json:"key,omitempty"`
}

type CountElements struct{
    NodeIP string   `json:"node,omitempty"`
    Count int     `json:"count,omitempty"`
}

func (hr HashCircle) Len() int {  
    return len(hr)  
}  
  
func (hr HashCircle) Less(i, j int) bool {  
    return hr[i] < hr[j]  
}  
  
func (hr HashCircle) Swap(i, j int) {  
    hr[i], hr[j] = hr[j], hr[i]  
}  
  
type Node struct {  
    Id       int  
    IP       string    
}  
  
func NewNode(id int, ip string) *Node {  
    return &Node{  
        Id:       id,  
        IP:       ip,  
    }  
}  
  
type ConsistentHash struct {  
    Nodes       map[uint32]Node  
    IsPresent   map[int]bool  
    Circle      HashCircle  
    
}  
  
func NewConsistentHash() *ConsistentHash {  
    return &ConsistentHash{  
        Nodes:     make(map[uint32]Node),   
        IsPresent: make(map[int]bool),  
        Circle:      HashCircle{},  
    }  
}  
  
func (hr *ConsistentHash) AddNode(node *Node) bool {  
 
    if _, ok := hr.IsPresent[node.Id]; ok {  
        return false  
    }  
    str := hr.ReturnNodeIP(node)  
    hr.Nodes[hr.GetHashValue(str)] = *(node)
    hr.IsPresent[node.Id] = true  
    hr.SortHashCircle()  
    return true  
}  
  
func (hr *ConsistentHash) SortHashCircle() {  
    hr.Circle = HashCircle{}  
    for k := range hr.Nodes {  
        hr.Circle = append(hr.Circle, k)  
    }  
    sort.Sort(hr.Circle)  
}  
  
func (hr *ConsistentHash) ReturnNodeIP(node *Node) string {  
    return node.IP 
}  
  
func (hr *ConsistentHash) GetHashValue(key string) uint32 {  
    return crc32.ChecksumIEEE([]byte(key))  
}  
  
func (hr *ConsistentHash) Get(key string) Node {  
    hash := hr.GetHashValue(key)  
    i := hr.SearchForNode(hash)  
    return hr.Nodes[hr.Circle[i]]  
}  

func (hr *ConsistentHash) SearchForNode(hash uint32) int {  
    i := sort.Search(len(hr.Circle), func(i int) bool {return hr.Circle[i] >= hash })  
    if i < len(hr.Circle) {  
        if i == len(hr.Circle)-1 {  
            return 0  
        } else {  
            return i  
        }  
    } else {  
        return len(hr.Circle) - 1  
    }  
}  
  
func PutKeyValues(circle *ConsistentHash){
        var str,input string
        fmt.Println("Enter the Key:")
        fmt.Scanf("%s\n",&str)
        fmt.Println("Enter the value:")
        fmt.Scanf("%s\n",&input)
        ipAddress := circle.Get(str)  
        address := "http://"+ipAddress.IP+"/keys/"+str+"/"+input
        req,err := http.NewRequest("PUT",address,nil)
        client := &http.Client{}
        resp, err := client.Do(req)
        if err!=nil{
            fmt.Println("Error:",err)
        }else{
            defer resp.Body.Close()
            fmt.Println(resp.StatusCode)
        }  
}  

func GetKeyValue(key string,circle *ConsistentHash){
    var out KeyValue 
    ipAddress:= circle.Get(key)
    response,err:= http.Get("http://"+ipAddress.IP+"/keys/"+key)
    if err!=nil{
        fmt.Println("Error:",err)
    }else{
        defer response.Body.Close()
        contents,err:= ioutil.ReadAll(response.Body)
        if(err!=nil){
            fmt.Println(err)
        }
        json.Unmarshal(contents,&out)
        result,_:= json.Marshal(out)
        fmt.Println(string(result))
    }
}

func GetAll(circle *ConsistentHash){
    var out []KeyValue
    response,err:= http.Get("http://127.0.0.1:3000/all")
    if err!=nil{
        fmt.Println("Error:",err)
    }else{
        defer response.Body.Close()
        contents,err:= ioutil.ReadAll(response.Body)
        if(err!=nil){
            fmt.Println(err)
        }
        json.Unmarshal(contents,&out)
        result,_:= json.Marshal(out)
        fmt.Println(string(result))
    }
}

func GetSpecificPort(){
     
    var port string
    fmt.Println("Enter the port no:")
    fmt.Scanf("%s\n",&port)
    var out []KeyValue
    response,err:= http.Get("http://127.0.0.1:"+port+"/keys")
    if err!=nil{
        fmt.Println("Error:",err)
    }else{
        defer response.Body.Close()
        contents,err:= ioutil.ReadAll(response.Body)
        if(err!=nil){
            fmt.Println(err)
        }
        json.Unmarshal(contents,&out)
        result,_:= json.Marshal(out)
        fmt.Println(string(result))
    }
}
func ShowCount(){
    var output[] CountElements
    response,err:= http.Get("http://127.0.0.1:3000/count")
    if err!=nil{
        fmt.Println("Error:",err)
    }else{
        defer response.Body.Close()
        contents,err:= ioutil.ReadAll(response.Body)
        if(err!=nil){
            fmt.Println(err)
        }
        json.Unmarshal(contents,&output)
        for _,j := range output{
            fmt.Println("\nNode/IP Address: "+j.NodeIP+" Count:"+strconv.Itoa(j.Count))
        } 
    }
}

func RemoveNode(circle *ConsistentHash){
    var port int
    fmt.Println("Enter the portno:")
    fmt.Scanf("%d\n",&port)
    var out []KeyValue
    response,err:= http.Get("http://127.0.0.1:"+strconv.Itoa(port)+"/remove")
    if err!=nil{
        fmt.Println("Error:",err)
    }else{
        defer response.Body.Close()
        contents,err:= ioutil.ReadAll(response.Body)
        if(err!=nil){
            fmt.Println(err)
        }
        json.Unmarshal(contents,&out)
    }
    for _,p:=range circle.Nodes{
        if(p.IP=="127.0.0.1:"+strconv.Itoa(port)){
            delete(circle.Nodes,circle.GetHashValue("127.0.0.1:"+strconv.Itoa(port)))
            delete(circle.IsPresent,p.Id) 
            circle.Circle = append(circle.Circle[:p.Id],circle.Circle[p.Id+1:]...)
            }
    }
    for _,j:=range out{
        str:= strconv.Itoa(j.Key)
        ipAddress := circle.Get(str)  
        address := "http://"+ipAddress.IP+"/keys/"+str+"/"+j.Value
        req,err := http.NewRequest("PUT",address,nil)
        client := &http.Client{}
        resp, err := client.Do(req)
        if err!=nil{
            fmt.Println("Error:",err)
        }else{
            defer resp.Body.Close()
        }
    }
    fmt.Println("Node was removed temporarily and Keys have been re-mapped.Check the re-mapping")
}

func ShowMapping(){
    var output[] Mapping
    response,err:= http.Get("http://127.0.0.1:3000/mapping")
    if err!=nil{
        fmt.Println("Error:",err)
    }else{
        defer response.Body.Close()
        contents,err:= ioutil.ReadAll(response.Body)
        if(err!=nil){
            fmt.Println(err)
        }
        json.Unmarshal(contents,&output)
        for _,j := range output{
            fmt.Println("\nKey "+ strconv.Itoa(j.Key)+" is mapped to "+j.NodeAddress)
        } 
    }   
}

func main() {   
    circle := NewConsistentHash()      
    circle.AddNode(NewNode(0, "127.0.0.1:3000"))
	circle.AddNode(NewNode(1, "127.0.0.1:3001"))
	circle.AddNode(NewNode(2, "127.0.0.1:3002")) 
    var cont string ="Y"
    for cont=="Y"{
    var choice int
    fmt.Println("1.PUT values\n2.GET a keyValue pair\n3.GET all KeyValues\n4.GET port specific keyValues\n5.Show Node Count\n6.Show Mapping\n7.Remove a Node")
    fmt.Println("Enter your choice:")
    fmt.Scanf("%d\n",&choice)
    switch choice{
        case 1:
            PutKeyValues(circle)
            break
        case 2:
            fmt.Println("\nEnter the Key:")
            var key string
            fmt.Scanf("%s\n",&key) 
            GetKeyValue(key,circle)
            break
        case 3:
            GetAll(circle)
            break
        case 4:
            GetSpecificPort()
            break
        case 5:
            ShowCount()
            break
        case 6:
            ShowMapping()
            break
        case 7:
            RemoveNode(circle)
            break
        default:
            break
        } 
    
    fmt.Println("Do you want to continue?(Y/N)")
    fmt.Scanf("%s\n",&cont)
    }
}  
