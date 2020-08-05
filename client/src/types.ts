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
  maximumFIATBuyAmount: number,
  minimumProfitPerSold: number,
  minimumPriceDropToBuy: number,
  growthIncreaseLimit: number,
  growthDecreaseLimit: number
}

export interface StatisticsOptions {
  numberOfPointsHold: number
}

export interface CollectorOptions {
  priceVariationDetection: number
  newPriceTimeRate: number
  datasource: null | string
}

export interface NotificationOptions {
  receiver: string,
  sender: string,
  senderpassword: string,
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

export interface ApplicationOptions {
  notificationOptions: NotificationOptions,
  statisticsOptions: StatisticsOptions,
  decisionMakerOptions: DecisionMakerOptions,
  collectorOptions: CollectorOptions,
}

export interface Application {
  _id: string,
  asset: string,
  accountID: string,
  options: ApplicationOptions
  createdAt: Date
}

export interface Account {
  _id: string,
  amount: Number,
  broker: string,
}

export interface ApplicationExecutionState {
  '_id': string,
  'executionId': string,
  'date': Date,
  'state': {
    'average': number,
    'standardDeviation': number,
    'lowerBollingerBand': number,
    'higherBollingerBand': number,
    'currentPrice': number,
    'currentChange': number
  }
}

export interface LogEvent {
  '_id': string,
  eventName: string,
  message: string,
  notified: boolean,
  createdAt: Date,
  applicatioID: string
}
