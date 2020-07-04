import axios from 'axios';
import React, { useEffect, useState } from 'react';

import {
  Benchmark, DataSourceOptions,
} from '../types';
import BenchmarkForm from './BenchmarkForm/BenchmarkForm';
import BenchmarkResult from './BenchmarkResult';
import BenchmarkList from './BenchmarksList';

export default () => {
  const [benchmarkResult, setBenchmarkResult] = useState<Benchmark>();

  const [dataSourceOptions, setDataSourceOptions] = useState<DataSourceOptions>();
  const [benchmarks, setBenchmarks] = useState<any>([]);


  useEffect(() => {
    fetch('/api/benchmark/data-sources')
      .then((res) => res.json())
      .then(setDataSourceOptions);

    fetch('/api/benchmark')
      .then((res) => res.json())
      .then((benchmarks) => setBenchmarks(benchmarks || []));
  }, []);


  // useEffect(() => {
  //   let newData: BenchmarkResult | undefined;
  //   if (benchmarkResult) {
  //     newData = JSON.parse(JSON.stringify(benchmarkResult));
  //     let start;
  //     let end;

  //     if (startDate) {
  //       start = new Date(startDate).getTime();
  //     }

  //     if (endDate) {
  //       end = new Date(endDate).getTime();
  //     }

  //     if (newData) {
  //       newData = filterBenchmarkResultByTime(newData, start, end);
  //     }
  //   }
  //   setData(newData);
  // }, [startDate, endDate, benchmarkResult]);

  const executeBenchmark = (input: any) => {
    axios
      .post('/api/benchmark', input)
      .then((res) => setBenchmarks(benchmarks.concat(res.data)));
  };

  return (
    <div>
      <div className="mt-3 mb-5">
        {
          dataSourceOptions
          && <BenchmarkForm onSubmit={executeBenchmark} dataSourceOptions={dataSourceOptions} />
        }
      </div>
      <div className="mt-3 mb-5">
        <BenchmarkList
          benchmarks={benchmarks}
          selectBenchmark={(id: string) => {
            const benchmark = benchmarks.find((b: any) => b._id === id);

            if (benchmark) {
              setBenchmarkResult(JSON.parse(JSON.stringify(benchmark)));
            }
          }}
          deleteBenchmark={(id: string) => {
            axios.delete(`/api/benchmark/${id}`)
              .then(() => {
                setBenchmarks(benchmarks.filter((b: any) => b._id !== id));
              });
          }}
        />
      </div>
      {
        benchmarkResult
        && (
        <div className="mb-5">
          <BenchmarkResult benchmark={benchmarkResult} />
        </div>
        )
      }
    </div>
  );
};
