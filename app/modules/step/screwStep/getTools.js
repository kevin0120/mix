import type { tPoint } from './interface/typeDef';
import { getDevice } from '../../deviceManager/devices';

export function getTools(points: Array<tPoint>) {
  const toolSnSet = new Set(points.map(p => p.tightening_tool));
  const lostTool = [];
  const tools = [];
  [...toolSnSet].forEach(t => {
    const tool = getDevice(t);
    if (tool) {
      tools.push(tool);
    } else {
      lostTool.push(t);
    }
  });

  if (lostTool.length > 0) {
    throw new Error(`tools not found: ${String(lostTool.map(t => `${t}`))}`);
  }

  const unhealthyTools = tools.filter(t => !t.Healthz);
  if (unhealthyTools.length > 0) {
    throw new Error(
      `tool not connected: ${JSON.stringify(
        unhealthyTools.map(t => `${String(t.serialNumber)}`)
      )}`
    );
  }
  return tools;
}
