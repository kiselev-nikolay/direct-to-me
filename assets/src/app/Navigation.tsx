import React from 'react';

import { ExternalLinkIcon } from '@chakra-ui/icons';
import {
  Box,
  Center,
  Divider,
  Heading,
  Image,
  Link,
} from '@chakra-ui/react';

import pages from './pagesIndex';

interface NavLinkProps {
  title: string;
}
let NavLink = (props: NavLinkProps) => {
  return (<>
    <Link mx=".5rem">{props.title}</Link>
    <Divider my=".5rem" />
  </>);
};

interface NavigationProps {
  page: number;
  setPage: any;
}

let Navigation = (props: NavigationProps) => {
  let sp = (page: string) => {
    localStorage.setItem("page", page);
    props.setPage(pages.indexOf(page));
  };
  return (
    <Center>
      <Box m="2rem" borderWidth="1px" borderRadius="lg" overflow="hidden">
        <Image mx="25%" my="1rem" w="50%" src="static/logo.png" alt="" />
        <Box m="1rem" flex="1" textAlign="left">
          <Heading mx=".5rem" as="h4" size="md">Direct to me</Heading>
          <Divider my=".5rem" />
          <div onClick={() => { sp("index"); }}>
            <NavLink title="Statistics" />
          </div>
          <div onClick={() => { sp("listRedirects"); }}>
            <NavLink title="List redirects" />
          </div>
          <div onClick={() => { sp("newRedirect"); }}>
            <NavLink title="New redirect" />
          </div>
          <Link mx=".5rem" href="https://nikolai.works" isExternal>
            Copyright <ExternalLinkIcon mx="2px" mt="-3px" />
          </Link>
        </Box>
      </Box>
    </Center>);
};

export default Navigation;