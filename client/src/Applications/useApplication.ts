import { useEffect, useState } from 'react';

import {
  Account, Application, ApplicationExecutionState, Asset, LogEvent,
} from '../types';

export default (applications: Application[]) => {
  const [activeApplicationID, setActiveApplicationID] = useState<string>();
  const [activeApplication, setActiveApplication] = useState<Application>();
  const [activeAccount, setActiveAccount] = useState<Account>();
  const [lastApplicationState, setLastApplicationState] = useState<ApplicationExecutionState>();
  const [assets, setAssets] = useState<Asset[]>();
  const [logEvents, setLogEvents] = useState<LogEvent[]>();

  useEffect(() => {
    if (!activeApplication) {
      setActiveAccount(undefined);
      return;
    }

    fetch(`/api/accounts/${activeApplication.accountID}`)
      .then((res) => res.json())
      .then(setActiveAccount);


    fetch(`/api/applications/${activeApplication._id}/state/last`)
      .then((res) => res.json())
      .then((data) => ({
        ...data,
        state: data.state.reduce((final: any, pair: any) => {
          final[pair.Key] = pair.Value;
          return final;
        }, {}),
      }))
      .then(setLastApplicationState);

    fetch(`/api/accounts/${activeApplication.accountID}/assets`)
      .then((res) => res.json())
      .then((assets) => assets || [])
      .then(setAssets);

    fetch(`/api/applications/${activeApplication._id}/log-events`)
      .then((res) => res.json())
      .then(setLogEvents);
  }, [activeApplication]);

  useEffect(() => {
    if (!activeApplicationID) {
      setActiveApplication(undefined);
      return;
    }

    const application = applications.find((app) => app._id === activeApplicationID);

    setActiveApplication(application);
  }, [activeApplicationID, applications]);

  return {
    activeApplication,
    setActiveApplication: setActiveApplicationID,
    activeAccount,
    lastApplicationState,
    assets,
    logEvents,
  };
};
