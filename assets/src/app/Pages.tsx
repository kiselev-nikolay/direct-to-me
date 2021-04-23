import React from 'react';

import { Box } from '@chakra-ui/react';

import NewRedirectFrom from './NewRedirect';

interface PagesProps {
  page: number;
}

function Pages(props: PagesProps) {
  return (
    <Box m="2rem" ml="0" borderWidth="1px" borderRadius="lg" overflow="hidden">
      <Box m="1rem" flex="1" textAlign="left">
        {props.page === 0 && <h1>Statistics</h1>}
        {props.page === 1 && <h1>List</h1>}
        {props.page === 2 && <NewRedirectFrom />}
      </Box>
    </Box>
  );
}

export default Pages;