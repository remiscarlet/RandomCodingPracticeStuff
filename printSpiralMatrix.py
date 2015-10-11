import copy

#Task: Given a 2d matrix, print out a list of the elements in order
# if the traversal is defined as "following the left wall" of said matrix.

# Eg, given the matrix:
#[[1,2,3],
# [4,5,6],
# [7,8,9]]

# output should be:
# [1,2,3,6,9,8,7,4,5]

# Invalid inputs (eg, empty list, non-2d list, etc) should return an empty list
# Can assume if given a matrix, it will not be a jagged matrix.

matrix1 = [[1,2,3],
           [4,5,6],
           [7,8,9]]

matrix2 = [[ 1, 2, 3, 4, 5],
           [ 6, 7, 8, 9,10],
           [11,12,13,14,15],
           [16,17,18,19,29]]

matrix3 = [[ 1, 2, 3, 4, 5, 6, 7],
           [ 8, 9,10,11,12,13,14]]

matrix4 = []

matrix5 = [[]]
def returnSpiralMatrix(matrix):
    if len(matrix)==0 or len(matrix[0]) == 0:
        return []
    #mapMatrix will mark "visited" "nodes" (Consider the matrix as a graph)
    mapMatrix = copy.deepcopy(matrix)
    for i in xrange(len(mapMatrix)):
        for j in xrange(len(mapMatrix[0])):
            mapMatrix[i][j] = False

    directionSwapMap = {"right":"down",
                        "down":"left",
                        "left":"up",
                        "up":"right"}

    outputList = []

    dirCoordMap = {"right":(0,+1),
                   "left":(0,-1),
                   "up":(-1,0),
                   "down":(+1,0)}

    width = len(matrix[0])
    height = len(matrix)


    # Function will recursively "crawl" through the matrix utilizing
    # a copied matrix of same dimensions but with the elements
    # being Bools marking whether that "node" has been visited
    # already or not. Using this, it simply traverses 
    # along the "graph" following the left "wall".

    # input d is a string from ["left","right","up","down"]
    # pos is a tuple of coordinates with (row,col)
    def recursiveCrawl(d,pos):
        row = pos[0]
        col = pos[1]
        mapMatrix[row][col] = True
        outputList.append(matrix[row][col])
        coordDelta = dirCoordMap[d]
        newRow = row+coordDelta[0]
        newCol = col+coordDelta[1]
        if newRow >= height:
            d = "left"
        elif newRow < 0:
            d = "right"
        elif newCol >= width:
            d = "down"
        elif newCol < 0:
            d = "up"

        coordDelta = dirCoordMap[d]
        newRow = row+coordDelta[0]
        newCol = col+coordDelta[1]

        if mapMatrix[newRow][newCol] == True:
            d = directionSwapMap[d]
            coordDelta = dirCoordMap[d]
            newRow = row+coordDelta[0]
            newCol = col+coordDelta[1]

        if mapMatrix[newRow][newCol] == True:
            return None
        else:
            return recursiveCrawl(d,(newRow,newCol))


    recursiveCrawl("right",(0,0))
    return outputList


print "Testing returnSpiralMatrix"
assert(returnSpiralMatrix(matrix1)==[1, 2, 3, 6, 9, 8, 7, 4, 5])
assert(returnSpiralMatrix(matrix2)==[1, 2, 3, 4, 5, 10, 15, 29, 19, 18, 17, 16, 11, 6, 7, 8, 9, 14, 13, 12])
assert(returnSpiralMatrix(matrix3)==[1, 2, 3, 4, 5, 6, 7, 14, 13, 12, 11, 10, 9, 8])
assert(returnSpiralMatrix(matrix4)==[])
assert(returnSpiralMatrix(matrix5)==[])
print "Tests passed!"
