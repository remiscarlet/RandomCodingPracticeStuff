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
    return &Node{
        is_root: is_root, 
        score: 0, 
        edges: make(map[byte]*Node, 0),
    }
}

func log(msg string) {

    IS_SILENT := true

    if ! IS_SILENT {
        fmt.Print(msg)
    }
}

func createSuffixLinks(curr_node *Node, root_node *Node) {
    for _, next_node := range curr_node.edges {
        log(fmt.Sprintf("Creating suffix links from node <%s>\n", next_node.node_name))
        var suffix []byte = next_node.node_name[1:]
        len_suffix := len(suffix)

        if len_suffix == 0 {
            log(fmt.Sprintf(">No suffix at <%s>. Linking to root.\n", next_node.node_name))
            next_node.suffix = root_node
        } else {
            for i := 0; i < len_suffix; i++ {
                log(fmt.Sprintf("Searching for sub-suffix <%s> in score trie\n", string(suffix)))
                target_node := searchForNode(string(suffix), root_node) 
                log(fmt.Sprintf(">Linking curr_node <%s> to target node <%s>\n", next_node.node_name, target_node.node_name))
                next_node.suffix = target_node
                suffix = suffix[1:]
                if target_node != root_node {
                    // Since we search suffixes from longest to shortest, if we ever get a target_node that isn't a root_node,
                    // it's a strictly longer suffix than all the ones we _haven't_ checked, hence skip them.
                    break
                }
            }
        }
        createSuffixLinks(next_node, root_node)
    }
}

func createDictSuffixes(curr_node *Node, root_node *Node) {
    for _, next_node := range curr_node.edges {
        log(fmt.Sprintf("@Attempting dict_suffix on node <%s>\n", next_node.node_name))
        connectDictSuffix(next_node, next_node)
        createDictSuffixes(next_node, root_node)
    }
}

func connectDictSuffix(curr_node *Node, orig_node *Node) {
    suffix := curr_node.suffix
    log(fmt.Sprintf("Suffix node <%s> is root: %t\n", suffix.node_name, suffix.is_root))
    if suffix != nil && ! suffix.is_root {
        if suffix.score == 0 {
            // If score is 0, this node is not in dict. Keep searching by traversing to suffix's suffix.
            connectDictSuffix(suffix, orig_node)
        } else {
            log(fmt.Sprintf("@#Creating dict suffix to node <%s> on node <%s>\n", suffix.node_name, orig_node.node_name))
            orig_node.dict_suffix = suffix
        }
    }
}

var NODE_SEARCH_MEMO = make(map[string]*Node)
func searchForNode(target_node_name string, root_node *Node) *Node {
    if node_ptr, ok := NODE_SEARCH_MEMO[target_node_name]; ok {
        return node_ptr
    } 

    curr_node := root_node

    var char byte
    var ok bool

    for i, _ := range target_node_name {
        char = target_node_name[i]
        if curr_node, ok = curr_node.edges[char]; ! ok {
            NODE_SEARCH_MEMO[target_node_name] = root_node
            return root_node
        }
    }

    NODE_SEARCH_MEMO[target_node_name] = curr_node

    return curr_node
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
    log(strings.Repeat("# # ", depth))
    log(fmt.Sprintf(">addr: %p, root: %t, score: %d, name: %s\n", node, node.is_root, node.score, string(node.node_name)))
    log(strings.Repeat("# # ", depth))
    log(fmt.Sprintf("+suffix: %p dict_suffix: %p \n", node.suffix, node.dict_suffix))
    for _, child_node := range node.edges {
        printNode(child_node, depth + 1)
    }
}

/**
 * TODO: Memoize this shit somehow.
 */
func generateScoreTrie(genes []string, health []int32) *Node {
    var root_node *Node = genNewNode(true)

    for i, gene := range genes {
        health_score := health[i]
        log(fmt.Sprintf("Adding gene(%s) with score: %d\n", gene, health_score))
        if (gene == "") {
            log(fmt.Sprintf("Empty gene???\n"))
        }
        addGeneToTrie(gene, health_score, root_node)
    }

    printTrie(root_node)

    return root_node
}

func scoreDNA(dna string, score_trie *Node, score_trie_root_node *Node) int32 {
    if dna == "" {
        return 0
    }

    var score int32 = 0

    var curr_node *Node = score_trie

    for i := 0; i < len(dna); i++ {
        log("--")
        var char byte = dna[i]
        if curr_node.is_root {
            log(fmt.Sprintf("Starting at <root> node. "))
        } else {
            log(fmt.Sprintf("Starting at node <%s>. ", curr_node.node_name))
        }
        log(fmt.Sprintf("Looking for transition function ->[%s]\n", string(char)))

        if next_node, ok := curr_node.edges[char]; ok {
            curr_node = next_node
            log(fmt.Sprintf("Now scoring for node %s on transition %d:%s\n", curr_node.node_name, i, string(char)))
            score += aggDictScore(curr_node)
        } else {
            if ! curr_node.is_root {
                // Root doesn't have a suffix. 
                curr_node = curr_node.suffix
                log(fmt.Sprintf("Transition ->[%s] did not exist in edges. Jump suffix link to node <%s>. Decrement pos\n", string(char), curr_node.node_name))
                i -= 1
            }
        }
    }

    return score
}

var AGG_DICT_SCORE_MEMO = make(map[*Node]int32)

func aggDictScore(start_node *Node) int32 {
    if score, ok := AGG_DICT_SCORE_MEMO[start_node]; ok {
        return score
    }

    if (start_node == nil) {
        log(fmt.Sprintf("Really shouldn't be getting nil passed in here...\n"))
    } else if ( ! start_node.is_root) {
        log(fmt.Sprintf("- On node %s with score %d\n", start_node.node_name, start_node.score))
    }
    
    var score int32
    if start_node.dict_suffix == nil {
        score = start_node.score
    } else {
        score = start_node.score + aggDictScore(start_node.dict_suffix)
    }
    AGG_DICT_SCORE_MEMO[start_node] = score
    return score
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

    log(fmt.Sprintf("%#v\n", genes))
    log(fmt.Sprintf("%#v\n", health))

    for sItr := 0; sItr < int(s); sItr++ {
        firstLastd := strings.Split(readLine(reader), " ")

        firstTemp, err := strconv.ParseInt(firstLastd[0], 10, 64)
        checkError(err)
        first := int32(firstTemp)

        lastTemp, err := strconv.ParseInt(firstLastd[1], 10, 64)
        checkError(err)
        last := int32(lastTemp)

        d := firstLastd[2]

        NODE_SEARCH_MEMO = make(map[string]*Node)
        AGG_DICT_SCORE_MEMO = make(map[*Node]int32)
        temp_genes := genes[first:last+1]
        temp_health := health[first:last+1]
        trie_root_node := generateScoreTrie(temp_genes, temp_health)
        createSuffixLinks(trie_root_node, trie_root_node)
        createDictSuffixes(trie_root_node, trie_root_node)

printTrie(trie_root_node)

        score := scoreDNA(d, trie_root_node, trie_root_node)

        log(fmt.Sprintf(">>>> The DNA strand %s has a score of %d\n", d, score))

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

