import React from 'react';
import { observer } from 'mobx-react-lite';
import { Outlet, Route, Routes } from 'react-router-dom';
import MainPage from './content/MainPage';


const App = observer(() => {
  return (
    <>
      <Routes>
        <Route path="/" element={<Outlet />}>
          <Route index element={<MainPage />} />
        </Route>
        <Route path="*" element={<div>Page not found</div>} />
      </Routes>
    </>
  );
});

export default App;