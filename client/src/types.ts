export interface Asset {
  _id: string
  amount: number
  buyTime: Date
  sellTime: Date
  buyPrice: number
  sellPrice: number
  sold: boolean
}

export interface BenchmarkOutput {
  finalAmount: number,
  balances: [number, number][],
  buys: [number, number][],
  sells: [number, number][],
  assets: Asset[]
}

export interface Benchmark {
  _id: string
  input: BenchmarkInput
  output: BenchmarkOutput
  status: string
  createdAt: Date
}

export interface DecisionMakerOptions {
  maximumBuyAmount: number,
  minimumProfitPerSold: number,
  minimumPriceDropToBuy: number,
  minutesToCollectNewPoint: number,
  growthIncreaseLimit: number,
  growthDecreaseLimit: number
}

export interface StatisticsOptions {
  numberOfPointsHold: number
}

export interface CollectorOptions {
  priceVariationDetection: number
}

export interface BenchmarkInput {
  decisionMakerOptions: DecisionMakerOptions,
  statisticsOptions: StatisticsOptions
  collectorOptions: CollectorOptions
  accountInitialAmount: number,
  dataSourceFilePath: string,
  asset: string,
}

export interface SelectOption {
  label: string,
  value: string,
}

export interface DataSourceOptions {
  [asset: string]: {
    [dataSource: string]: string
  }
}

export interface DatesInterval {
  startDate: Date,
  endDate: Date,
}
