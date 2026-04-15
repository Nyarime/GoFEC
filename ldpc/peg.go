package ldpc

// PEG (Progressive Edge Growth) LDPC矩阵构造
// 贪婪最大化girth，提升BP解码收敛性能

// BuildPEGMatrix 构造PEG优化的稀疏校验矩阵
// n=变量节点数(码长), m=校验节点数, dv=每变量节点度
func BuildPEGMatrix(n, m, dv int) [][]bool {
	// 邻接表
	varToCheck := make([][]int, n) // variable → connected checks
	checkToVar := make([][]int, m) // check → connected variables

	for v := 0; v < n; v++ {
		for d := 0; d < dv; d++ {
			// 找girth最大的check节点
			bestC := -1
			bestDepth := -1

			for c := 0; c < m; c++ {
				// 跳过已连接的
				if connected(varToCheck[v], c) { continue }
				// 限制check度
				if len(checkToVar[c]) >= (n*dv/m + 2) { continue }

				depth := bfsDepth(v, c, varToCheck, checkToVar, n, m)
				if depth > bestDepth {
					bestDepth = depth
					bestC = c
				}
			}

			if bestC == -1 {
				// 所有check都满或已连，选度最小的
				minDeg := n
				for c := 0; c < m; c++ {
					if !connected(varToCheck[v], c) && len(checkToVar[c]) < minDeg {
						minDeg = len(checkToVar[c])
						bestC = c
					}
				}
			}
			if bestC == -1 { break }

			varToCheck[v] = append(varToCheck[v], bestC)
			checkToVar[bestC] = append(checkToVar[bestC], v)
		}
	}

	// 转换为矩阵
	matrix := make([][]bool, m)
	for c := 0; c < m; c++ {
		matrix[c] = make([]bool, n+m)
		for _, v := range checkToVar[c] {
			matrix[c][v] = true
		}
		matrix[c][n+c] = true // 对角线校验位
	}
	return matrix
}

// bfsDepth BFS计算v到c的最短路径深度(近似girth)
func bfsDepth(v, c int, v2c [][]int, c2v [][]int, n, m int) int {
	if len(v2c[v]) == 0 { return 100 } // 未连接=无穷girth

	visited := make(map[int]bool)
	queue := []int{v}
	visited[v] = true
	depth := 0

	for len(queue) > 0 && depth < 10 { // 最多搜10层
		depth++
		next := []int{}
		for _, vv := range queue {
			for _, cc := range v2c[vv] {
				if cc == c { return depth }
				for _, vvv := range c2v[cc] {
					if !visited[vvv] {
						visited[vvv] = true
						next = append(next, vvv)
					}
				}
			}
		}
		queue = next
	}
	return 100 // 未到达=大girth
}

func connected(list []int, target int) bool {
	for _, v := range list {
		if v == target { return true }
	}
	return false
}
