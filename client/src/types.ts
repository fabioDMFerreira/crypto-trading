export interface Asset {
  _id: string
  amount: number
  buyTime: Date
  sellTime: Date
  buyPrice: number
  sellPrice: number
  sold: boolean
}

export interface BenchmarkResult {
  prices: [number, number][],
  balances: [number, number][],
  buys: [number, number][],
  sells: [number, number][],
  assets: Asset[]
}

export interface BenchmarkInput {
  accountInitialAmount: number
  dataSourceFilePath: string
}

export interface BenchmarkOutput {
  finalAmount: number
}

export interface Benchmark {
  _id: string
  input: BenchmarkInput
  output: BenchmarkOutput
  status: string
  createdAt: Date
}
