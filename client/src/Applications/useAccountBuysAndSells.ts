import { useEffect, useState } from 'react';

export default (accountID: string) => {
  const [buys, setBuys] = useState<[number, number][]>([]);
  const [sells, setSells] = useState<[number, number][]>([]);


  useEffect(() => {
    if (accountID) {
      fetch(`/api/accounts/${accountID}/buys-and-sells`)
        .then((res) => res.json())
        .then((data: any) => {
          setBuys(data.buys);
          setSells(data.sells);
        });
    } else {
      setBuys([]);
      setSells([]);
    }
  }, [accountID]);

  return {
    buys,
    sells,
  };
};
