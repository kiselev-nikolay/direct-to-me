import React from 'react';

import Chart from 'react-apexcharts';

import {
  Box,
  Divider,
  Flex,
  Heading,
  Link,
  SimpleGrid,
  Skeleton,
  Stack,
  Table,
  Tag,
  Tbody,
  Td,
  Tr,
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

interface ApexConf {
  legend?: {
    show: boolean;
  },
  dataLabels?: {
    enabled: boolean,
  };
  series: Array<number>,
  labels?: Array<string>;
}
type ApexData = Array<{
  key: string,
  clicksChart: ApexConf,
  clicks: APIClick,
  fails: APIFail;
}>;

function makeCharts(inputData: Map<string, APIStats>) {
  let data: ApexData = [];
  let clicks: ApexConf = {
    dataLabels: { enabled: true },
    legend: { show: true },
    series: [0, 0],
    labels: ["Direct", "Social"]
  };
  let redirects: ApexConf = {
    dataLabels: { enabled: true },
    legend: { show: true },
    series: [],
    labels: []
  };
  for (let [k, v] of inputData) {
    let dataRow: any = { key: k };
    if (v.clicks !== null && v.clicks !== undefined) {
      // Add data to direct/social global compare chart
      clicks.series[0] += v.clicks.direct;
      clicks.series[1] += v.clicks.social;
      // Add data to redirects global compare chart
      redirects.series.push(v.clicks.direct + v.clicks.social);
      redirects.labels.push(k);
      // Add data for small individual redirect chart
      let clicksChart: ApexConf = { dataLabels: { enabled: false }, legend: { show: false }, series: [], labels: [] };
      if (v.clicks.direct != 0) {
        clicksChart.series.push(v.clicks.direct);
        clicksChart.labels.push("Direct");
      }
      if (v.clicks.social != 0) {
        clicksChart.series.push(v.clicks.social);
        clicksChart.labels.push("Social");
      }
      dataRow.clicksChart = clicksChart;
    }
    dataRow.clicks = v.clicks || {
      direct: 0,
      social: 0,
    };
    dataRow.fails = v.fails || {
      notFound: 0,
      databaseUnreachable: 0,
      templateProcessFailed: 0,
      clientContentProcessFailed: 0,
    };
    data.push(dataRow);
  }
  return { data: data, redirects: redirects, clicks: clicks };
}

interface StatLineProps {
  title: string;
  value: number;
}
function OkStatLine(props: StatLineProps) {
  return (<Tr>
    <Td>{props.title}</Td>
    <Td isNumeric><Tag>{props.value}</Tag></Td>
  </Tr>);
}
function ErrStatLine(props: StatLineProps) {
  if (props.value === 0) {
    return (<></>);
  }
  return (<Tr>
    <Td>{props.title}</Td>
    <Td isNumeric><Tag colorScheme="red">{props.value}</Tag></Td>
  </Tr>);
}

interface StatsProps {
  data: Map<string, APIStats>;
}

function Stats(props: StatsProps) {
  let { data, redirects, clicks } = makeCharts(props.data);
  return (<Stack>
    <style>{".apexcharts-legend.apexcharts-align-center.position-right{width:25%}"}</style>
    <SimpleGrid columns={2} spacing={3}>
      <Box flex="1" p="1rem" overflow="hidden">
        <Heading as="h2" size="md">Redirect total</Heading>
        <Divider my={3} />
        <Chart options={redirects} series={redirects.series} type="donut" />
      </Box>
      <Box flex="1" p="1rem" overflow="hidden">
        <Heading as="h2" size="md">Source total</Heading>
        <Divider my={3} />
        <Chart options={clicks} series={clicks.series} type="donut" />
      </Box>
    </SimpleGrid>
    <Box pb="3rem">
    </Box>
    {data.map((i) => {
      return (
        <Flex key={i.key} borderWidth="1px" borderRadius="lg" p="1rem" overflow="hidden">
          <Box flex="3">
            <Heading as="h2" size="md">Route <Link href={"/" + i.key} target="_blank">/{i.key}</Link></Heading>
            <Divider mt={3} />
            <Table variant="simple">
              <Tbody>
                <OkStatLine title="Direct redirects" value={i.clicks.direct} />
                <OkStatLine title="Social networks redirects" value={i.clicks.direct} />
                <ErrStatLine title="Not Found" value={i.fails.notFound} />
                <ErrStatLine title="Database Unreachable" value={i.fails.databaseUnreachable} />
                <ErrStatLine title="Template Process Failed" value={i.fails.templateProcessFailed} />
                <ErrStatLine title="Client Content Process Failed" value={i.fails.clientContentProcessFailed} />
              </Tbody>
            </Table>
          </Box>
          {i.clicksChart && <Box flex="1">
            <Chart options={i.clicksChart} series={i.clicksChart.series} type="donut" />
          </Box>}
        </Flex>
      );
    })}
  </Stack>);
}