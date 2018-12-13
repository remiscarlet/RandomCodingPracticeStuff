package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
    "math"
    //"encoding/json"
)

type Node struct {
    is_root     bool
    score       int32
    node_name   []byte
    dict_suffix *Node
    suffix      *Node
    edges       map[byte]*Node
}

func genNewNode(is_root bool) *Node {
    return &Node{is_root: is_root, score: 0, edges: make(map[byte]*Node)}
}

func createSuffixLinks(curr_node *Node, root_node *Node) {
    for _, next_node := range curr_node.edges {
        fmt.Printf("Creating suffix links from node %s\n", next_node.node_name)
        var suffix []byte = next_node.node_name[1:]
        len_suffix := len(suffix)

        if len_suffix == 0 {
            fmt.Printf("No suffix at %s. Linking to root.\n", next_node.node_name)
            next_node.suffix = root_node
        } else {
            fmt.Printf("Searching for suffix %s\n", string(suffix))
            for i := 0; i < len_suffix; i++ {
                fmt.Printf("Searhing for subsuffix %s in score trie\n", string(suffix))
                target_node := searchForNode(string(suffix), root_node)
                if target_node != nil {
                    next_node.suffix = target_node
                } else {
                    fmt.Printf("Linking curr node %s to root.\n", next_node.node_name)
                    next_node.suffix = root_node
                }

                suffix = suffix[1:]
            }
        }

        createSuffixLinks(next_node, root_node)
    }
}

func createDictSuffixes(curr_node *Node, root_node *Node) {
    for _, next_node := range curr_node.edges {
        if ! next_node.is_root {
            connectDictSuffix(curr_node, curr_node)
        }
        createDictSuffixes(next_node, root_node)
    }
}

func connectDictSuffix(curr_node *Node, orig_node *Node) {
    suffix := curr_node.suffix
    if suffix != nil && ! suffix.is_root {
        orig_node.dict_suffix = suffix
        connectDictSuffix(suffix, orig_node)
    }
}

func searchForNode(target_node_name string, curr_node *Node, ) *Node {
    if string(curr_node.node_name) == target_node_name {
        return curr_node
    }

    curr_index := (len(target_node_name) - 1) - len(curr_node.node_name)

    var target_node *Node
    for char, child_node := range curr_node.edges {
        if target_node_name[curr_index] == char {
            target_node = searchForNode(target_node_name, child_node)
            if target_node != nil {
                break
            } 
        }
    }
    return target_node
}

func addGeneToTrie(gene string, score int32, root_node *Node) {
    addremGeneToTrie(gene, make([]byte, 0), score, root_node, true)
}

func removeGeneFromTrie(gene string, score int32, root_node *Node) {
    addremGeneToTrie(gene, make([]byte, 0), score, root_node, false)
}

func addremGeneToTrie(gene_remaining string, gene_removed []byte, score int32, parent_node *Node, isAddMode bool) {
    var first_char byte = gene_remaining[0]
    
    curr_node, ok := parent_node.edges[first_char]

    gene_removed = append(gene_removed, first_char)

    if ! ok {
        curr_node = genNewNode(false)
        parent_node.edges[first_char] = curr_node
        curr_node.node_name = gene_removed
    }

    if len(gene_remaining) == 1 {
        if isAddMode {
            curr_node.score += score
        } else {
            curr_node.score -= score
        }
    } else {
        addremGeneToTrie(gene_remaining[1:], gene_removed, score, curr_node, isAddMode)
    }
}

func printTrie(root_node *Node) {
    printNode(root_node, 0)
}

func printNode(node *Node, depth int) {
    fmt.Print(strings.Repeat("--  ", depth))
    fmt.Printf("addr: %p, root: %t, score: %d, name: %s\n", node, node.is_root, node.score, string(node.node_name))
    fmt.Print(strings.Repeat("--  ", depth))
    fmt.Printf("suffix: %p dict_suffix: %p \n", node.suffix, node.dict_suffix)
    for _, child_node := range node.edges {
        printNode(child_node, depth + 1)
    }
}

func generateScoreTrie(genes []string, health []int32) *Node {
    var root_node *Node = genNewNode(true)

    for i, gene := range genes {
        health_score := health[i]
        addGeneToTrie(gene, health_score, root_node)

        //fmt.Sprintf("Added gene(%s) with score: %d\n", gene, health_score)
    }

    //printTrie(root_node)

    return root_node
}
func scoreDNA2(dna string, score_trie *Node) int32 {
    if dna == "" {
        return 0
    }

    var score int32 = 0

    var curr_node *Node = score_trie
    var ok bool

    for i := 0; i < len(dna); i++ {
        available_edges := make([]byte, 0)
        
        if curr_node != nil {
            for char, _ := range curr_node.edges {
                available_edges = append(available_edges, char)
            }
        } else {
            fmt.Printf("What. Curr_node was nil.\n")
        }

        var char byte = dna[i]
        fmt.Printf("On transition step: %s. Available edges: %v\n", string(char), available_edges)

        if len(available_edges) == 0 {
            if curr_node.dict_suffix != nil {
                curr_node = curr_node.dict_suffix
                score += curr_node.score
            } else {
                curr_node = curr_node.suffix
            }
            for char, _ := range curr_node.edges {
                available_edges = append(available_edges, char)
            }
        }
        
        
        curr_node, ok = curr_node.edges[char]
        
        if ok {
            fmt.Printf("Adding score of %d at node %s\n", curr_node.score, curr_node.node_name)
            score += curr_node.score
        } else if curr_node != nil {
            fmt.Printf("Moving node to dict_suffix node %s\n", curr_node.dict_suffix.node_name)
            curr_node = curr_node.dict_suffix
            score += curr_node.score
            i -= 1
        }
    }

    return score
}

func scoreDNA(dna string, score_trie *Node) int32 {
    if dna == "" {
        return 0
    }

    var score int32 = 0

    var curr_node *Node = score_trie

    for i := 0; i < len(dna); i++ {
        var char byte = dna[i]

        if next_node, ok := curr_node.edges[char]; ok {
            curr_node = next_node
        } else {
            if ! curr_node.is_root {
                curr_node = curr_node.suffix
                i -= 1
            }
        }
        fmt.Println("--")
        score += aggDictScore(curr_node)
    }

    return score
}

func aggDictScore(start_node *Node) int32 {
    if (start_node == nil) {
    } else {
        fmt.Printf("On dict suffix node %s of score %d\n", start_node.node_name, start_node.score)
    }
    if start_node.dict_suffix == nil {
        return start_node.score
    } else {
        return start_node.score + aggDictScore(start_node.dict_suffix)
    }
}

func main() {
    reader := bufio.NewReaderSize(os.Stdin, 1024 * 1024)

    nTemp, err := strconv.ParseInt(readLine(reader), 10, 64)
    checkError(err)
    n := int32(nTemp)

    genesTemp := strings.Split(readLine(reader), " ")

    var genes []string

    for i := 0; i < int(n); i++ {
        genesItem := genesTemp[i]
        genes = append(genes, genesItem)
    }

    healthTemp := strings.Split(readLine(reader), " ")

    var health []int32

    for i := 0; i < int(n); i++ {
        healthItemTemp, err := strconv.ParseInt(healthTemp[i], 10, 64)
        checkError(err)
        healthItem := int32(healthItemTemp)
        health = append(health, healthItem)
    }

    sTemp, err := strconv.ParseInt(readLine(reader), 10, 64)
    checkError(err)
    s := int32(sTemp)

    var min_score int32 = math.MaxInt32
    var max_score int32 = math.MinInt32

    for sItr := 0; sItr < int(s); sItr++ {
        firstLastd := strings.Split(readLine(reader), " ")

        firstTemp, err := strconv.ParseInt(firstLastd[0], 10, 64)
        checkError(err)
        first := int32(firstTemp)

        lastTemp, err := strconv.ParseInt(firstLastd[1], 10, 64)
        checkError(err)
        last := int32(lastTemp)

        d := firstLastd[2]

        temp_genes := genes[first:last+1]
        temp_health := health[first:last+1]
        trie_root_node := generateScoreTrie(temp_genes, temp_health)
        createSuffixLinks(trie_root_node, trie_root_node)
        createDictSuffixes(trie_root_node, trie_root_node)

printTrie(trie_root_node)

        score := scoreDNA(d, trie_root_node)

        if min_score > score {
            min_score = score
        }
        if max_score < score {
            max_score = score
        }
    }
    fmt.Printf("%d %d\n", min_score, max_score)
}

func readLine(reader *bufio.Reader) string {
    str, _, err := reader.ReadLine()
    if err == io.EOF {
        return ""
    }

    return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
    if err != nil {
        panic(err)
    }
}
