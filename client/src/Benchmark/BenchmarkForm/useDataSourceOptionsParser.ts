import { useEffect, useState } from 'react';

import { DataSourceOptions, SelectOption } from '../../types';

export default (dataSourceOptions: DataSourceOptions) => {
  const [assets, setAssets] = useState<SelectOption[]>([]);
  const [dataSources, setDataSources] = useState<SelectOption[]>([]);
  const [activeAsset, setActiveAsset] = useState<string>();
  const [activeDataSource, setActiveDataSource] = useState<string>();


  useEffect(() => {
    if (dataSourceOptions) {
      const assets = Object.keys(dataSourceOptions);

      const assetsOptions: SelectOption[] = assets.map((asset) => ({ label: asset, value: asset }));

      setAssets(assetsOptions);

      // const dataSources = assets.reduce(
      //   (allDataSources, asset) => {
      //     const assetDataSources: SelectOption[] = Object.keys(dataSourceOptions[asset]).map((dataSourceKey) => ({
      //       label: dataSourceKey,
      //       value: dataSourceOptions[asset][dataSourceKey],
      //     }));

      //     return allDataSources.concat(assetDataSources);
      //   },
      //   [] as SelectOption[],
      // );

      setActiveAsset('btc');
    }
  }, [dataSourceOptions]);

  useEffect(() => {
    if (activeAsset) {
      const dataSources = Object.keys(dataSourceOptions[activeAsset])
        .map(
          (dataSourceKey) => ({
            label: dataSourceKey,
            value: dataSourceOptions[activeAsset][dataSourceKey],
          }),
        );
      setDataSources(dataSources);
      if (dataSources.length) {
        setActiveDataSource(dataSources[0].value)
      }
    } else {
      setDataSources([]);
    }
  }, [activeAsset, dataSourceOptions]);

  return {
    assets, dataSources, activeAsset, setActiveAsset, activeDataSource, setActiveDataSource
  };
};
