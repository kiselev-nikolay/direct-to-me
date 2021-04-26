import React from 'react';

import {
  Skeleton,
  Stack,
} from '@chakra-ui/react';

import GetStats, {
  APIClick,
  APIFail,
  APIStats,
} from './API/Stats';

let stats: Map<string, APIStats> = new Map<string, APIStats>();

interface StatisticsProps { }
interface StatisticsState {
  stats: Map<string, APIStats>;
}

export default class Statistics extends React.Component<StatisticsProps, StatisticsState> {
  constructor(props: StatisticsProps) {
    super(props);
    this.state = { stats: stats };
  }
  componentDidMount() {
    if (stats.size === 0) {
      setTimeout(() => {
        GetStats().then((s: Map<string, APIStats>) => { stats = s; this.setState({ stats: s }); });
      }, 500);
    } else {
      this.setState({ stats: stats });
    }
  }
  render() {
    return (<>
      {stats.size === 0 &&
        <Stack>
          <Skeleton height="20px" />
          <Skeleton height="20px" />
          <Skeleton height="20px" />
        </Stack>
      }
      <Stack>
        <Stats data={this.state.stats} />
      </Stack>
    </>);
  }
}

interface StatsProps {
  data: Map<string, APIStats>;
}

function Stats(props: StatsProps) {
  let data: Array<{
    key: string, clicks: APIClick,
    fails: APIFail;
  }> = [];
  for (let [k, v] of props.data) {
    data.push({ key: k, ...v });
  }
  return (<>
    {JSON.stringify(data)}
  </>);
}