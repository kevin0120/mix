using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using SVW.COMM;

namespace plugin
{
    enum PMON_STATUS
    {
        EVT_SO = 0,
        EVT_AO,
        EVT_SC,
        EVT_AC,
        EVT_SD,
        EVT_AD,
        EVT_UNEXCEPTED,
        EVT_MSG_TIMEOUT,
        EVT_DATA_RCV,
        EVT_DATA_SEND_SUCCESS,
        EVT_DATA_SEND_FAIL,
        EVT_CHANNEL_OPENED,
        EVT_CHANNEL_CLOSED,
        EVT_ERROR
    }

    class FisClinet
    {
        // pmon事件响应
        public int OnPmonEvent(int evt, string channel, string content)
        {
            if (content == null)
            {
                content = "";
            }
            Console.WriteLine(evt.ToString() + "::" + channel.ToString() + "::" + content.ToString());
            return 0;
        }

        public void Init()
        {
            // 跟踪所有日志
            if (PMON.PmonInit(new PMON.EventHandler(OnPmonEvent), 0))
            {
                Console.WriteLine("init pmon ok");
            }
            else
            {
                Console.WriteLine("init pmon failed");
            }
        }
    }
}
