import { useEffect, useState } from 'react';
import { Application } from '../types';

export default () => {
  const [applications, setApplications] = useState<Application[]>([]);

  useEffect(() => {
    fetch('/api/applications')
      .then((res) => res.json())
      .then(setApplications);
  }, []);

  function deleteApplicationByID(id: string) {
    fetch(`/api/applications/${id}`, { method: "delete" })
      .then(() => {
        setApplications(applications.filter(app => app._id !== id))
      })
  }

  return {
    applications,
    deleteApplicationByID
  };
};
