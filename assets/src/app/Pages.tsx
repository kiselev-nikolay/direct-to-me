import React from 'react';

import { Box } from '@chakra-ui/react';

import ListRedirects from './ListRedirects';
import NewRedirectFrom from './NewRedirect';
import pages from './pagesIndex';
import Statistics from './Statistics';

interface PagesProps {
  page: number;
}

function Pages(props: PagesProps) {
  return (
    <Box m="2rem" ml="0" borderWidth="1px" borderRadius="lg" overflow="hidden">
      <Box m="1rem" flex="1" textAlign="left">
        {props.page === pages.indexOf("index") && <Statistics />}
        {props.page === pages.indexOf("listRedirects") && <ListRedirects />}
        {props.page === pages.indexOf("newRedirect") && <NewRedirectFrom />}
      </Box>
    </Box>
  );
}

export default Pages;