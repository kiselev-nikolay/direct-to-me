import React from 'react';

import {
  Grid,
  GridItem,
} from '@chakra-ui/react';

import Navigation from './Navigation';
import Pages from './Pages';

interface AppProps { }
interface AppState {
  page: number;
}
class App extends React.Component<AppProps, AppState> {
  constructor(props: AppProps) {
    super(props);
    this.state = { page: 0 };
  }
  setPage(page: number) {
    this.setState({ page: page });
    return;
  }
  render() {
    return (
      <Grid
        h="200px"
        templateRows="repeat(1, 1fr)"
        templateColumns="repeat(5, 1fr)"
        gap={4}
      >
        <GridItem colSpan={1}>
          <Navigation page={this.state.page} setPage={(n: number) => this.setPage(n)}></Navigation>
        </GridItem>
        <GridItem colSpan={4}>
          <Pages page={this.state.page} />
        </GridItem>
      </Grid >
    );
  }
}

export default App;