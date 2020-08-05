import React from 'react';

import ApplicationItem from './ApplicationItem';
import ApplicationsList from './ApplicationsList';
import useApplication from './useApplication';
import useApplicationsList from './useApplicationsList';

export default () => {
  const { applications, deleteApplicationByID } = useApplicationsList();
  const {
    activeApplication,
    setActiveApplication,
    activeAccount,
    lastApplicationState,
    assets,
    logEvents,
  } = useApplication(applications);

  return (
    <>
      <ApplicationsList
        applications={applications}
        selectApplication={setActiveApplication}
        deleteApplication={deleteApplicationByID}
      />
      {
        activeApplication && activeAccount && lastApplicationState && assets && logEvents
        && (
          <div className="mt-5">
            <ApplicationItem
              application={activeApplication}
              account={activeAccount}
              lastApplicationState={lastApplicationState}
              assets={assets}
              logEvents={logEvents}
            />
          </div>
        )
      }
    </>
  );
};
