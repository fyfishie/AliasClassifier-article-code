'''
Author: fyfishie
LastEditors: fyfishie
LastEditTime: 2023-06-28:17
Description: :)
email: fyfishie@outlook.com
'''
import matplotlib.pyplot as plt
import csv
import math
format = 'pdf'
dpi = 400
plt.rcParams['font.sans-serif']=['SimHei']

def editDistance():
    fig,ax=plt.subplots(figsize=(8,6))
    plt.clf()
    X = []
    Y = []
    with open('./distance.csv','r') as file:
        plots = csv.reader(file,delimiter=',')
        for row in plots:
            X.append(int(row[0]))
            Y.append(float(row[1]))   
    plt.plot(X,Y)
    # plt.axhline(1,c='y',linestyle='--')
    plt.xlabel('路径编辑距离',fontsize=16)
    plt.ylabel('累积分布函数',fontsize=16)
    for x,y in zip(X,Y):
        plt.scatter(x,y,c='b')
        # if x<3 and x!=1:
        #     plt.text(x+0.8,y-0.04,(x,y),ha='center',va='bottom',fontsize=10)
        # if x==3:
        #     plt.text(x+0.2,y-0.05,(x,y),ha='center',va='bottom',fontsize=10)
        if x==4:
            plt.text(x+0.2,y-0.06,(x,y),ha='center',va='bottom',fontsize=15)
    plt.tick_params(labelsize=13)
    plt.savefig("./pac/pic2."+format,bbox_inches='tight',format=format,dpi=dpi)


def factor():
    plt.clf()
    plt.rcParams['font.sans-serif']=['SimHei']
    PX = []
    PY = []
    UX = []
    UY = []
    with open('./pair/factor/london.csv','r') as file:
        plots = csv.reader(file,delimiter=',')
        for row in plots:
            PX.append(float(row[0]))
            PY.append(float(row[1]))   
    with open('./upair/factor/london.csv','r') as file:
        plots = csv.reader(file,delimiter = ',')
        for row in plots:
            UX.append(float(row[0]))
            UY.append(float(row[1]))
    for x,y in zip(UX,UY):
        if x==0.5:
            plt.scatter(x,y,c='r',s=8)
            plt.text(x+0.1,y-0.04,(x,y),ha='center',va='bottom',fontsize=12)
    for x,y in zip(PX,PY):
        if x==0.5:
            plt.scatter(x,y,c='b',s=8)
            plt.text(x+0.1,y-0.04,(x,y),ha='center',va='bottom',fontsize=12)
    # plt.axline([0,1],[0.5,1],color='y',linestyle='--')
    plt.plot(PX,PY,'bs-',label = '别名IP')
    plt.axline([0.5,0.8],[0.5,0.311],color='y',linestyle='--')
    plt.plot(UX,UY,'ro--',label = '非别名IP')
    # axes.set_xlabel('x', fontdict=xlabel_font, labelpad=20, loc='right')
    plt.xlabel('路径相似系数',fontsize=16)
    plt.ylabel('累积分布函数',fontsize=16)
    plt.tick_params(labelsize=13)
    plt.legend(loc='lower right',fontsize=13)
    plt.savefig("./pac/pic4."+format,bbox_inches='tight',format=format,dpi=dpi)
    
def rtt():
    plt.clf()
    plt.rcParams['font.sans-serif']=['SimHei']
    PX = []
    PY = []
    UX = []
    UY = []
    with open('./pair/rtt/london.csv','r') as file:
        plots = csv.reader(file,delimiter=',')
        for row in plots:
            if int(row[0])>200:
                break
            if (int(row[0])%10)!=0:
                continue
            PX.append(int(row[0]))
            PY.append(float(row[1]))   
    with open('./upair/rtt/london.csv','r') as file:
        plots = csv.reader(file,delimiter = ',')
        for row in plots:
            if int(row[0])>200:
                break
            if int(row[0])%10!=0:
                continue
            UX.append(int(row[0]))
            UY.append(float(row[1]))
    for x,y in zip(PX,PY):
        # if x%10!=0:
        #     continue
        if x==180:
            plt.scatter(x,y)
            plt.text(x-22,y-0.07,(x,y),ha='center',va='bottom',fontsize=12)
    for x,y in zip(UX,UY):
        # if x%10!=0:
        #     continue
        if x==180:
            plt.scatter(x,y)
            plt.text(x+11,y-0.07,(x,y),ha='center',va='bottom',fontsize=12)
    plt.axline([180,0.8],[180,0.311],color='y',linestyle='--')
    plt.plot(PX,PY,'bs-',label = '别名IP')
    plt.plot(UX,UY,'ro--',label = '非别名IP')
    plt.xlabel('RTT差值',fontsize=16)

    plt.ylabel('累积分布函数',fontsize=16)
    plt.legend(loc='best',fontsize=13)
    plt.tick_params(labelsize=13)
    plt.savefig("./pac/pic1."+format,bbox_inches='tight',format=format,dpi=dpi)
    
def ttl():
    plt.clf()
    plt.rcParams['font.sans-serif']=['SimHei']
    PX = []
    PY = []
    UX = []
    UY = []
    with open('./pair/ttl/london.csv','r') as file:
        plots = csv.reader(file,delimiter=',')
        for row in plots:
            if int(row[0])>180:
                if int(row[0])%10!=0:
                    continue
            if int(row[0])<=15:
                if int(row[0])%3!=0:
                    continue
            if int(row[0]) >15 and int(row[0])<180:
                if int(row[0]) % 20 != 0:
                    continue
            PX.append(int(row[0]))
            PY.append(float(row[1]))   
    with open('./upair/ttl/london.csv','r') as file:
        plots = csv.reader(file,delimiter = ',')
        for row in plots:
            if int(row[0]) > 180:
                if int(row[0]) % 10 != 0:
                    continue
            if int(row[0]) <=15:
                if int(row[0]) % 3 != 0:
                    continue
            if int(row[0]) >15 and int(row[0])<180:
                if int(row[0]) % 20 != 0:
                    continue
            UX.append(int(row[0]))
            UY.append(float(row[1]))
    for x,y in zip(PX,PY):
        if x==15:
            plt.scatter(x,y,c='b',s=4)
            plt.text(x+16,y-0.06,(x,y),ha='center',va='bottom',fontsize=12)
    for x,y in zip(UX,UY):
        if x==15:
            plt.scatter(x,y,c='r',s=4)
            plt.text(x+16,y-0.06,(x,y),ha='center',va='bottom',fontsize=12)
    plt.axline([15,0.8],[15,0.311],color='y',linestyle='--')
    plt.plot(PX,PY,'bs-',label = '别名IP')
    plt.plot(UX,UY,'ro--',label = '非别名IP')
    plt.xlabel('Reply TTL差值',fontsize=16)
    plt.ylabel('累积分布函数',fontsize=16)
    plt.legend(loc='best',fontsize=13)
    plt.tick_params(labelsize=13)
    plt.savefig("./pac/pic5."+format,bbox_inches='tight',format=format,dpi=dpi)

def domain():
    plt.clf()
    plt.rcParams['font.sans-serif']=['SimHei']
    PX = []
    PY = []
    UX = []
    UY = []
    with open('./pair/domain/distance.txt','r') as file:
        plots = csv.reader(file,delimiter=',')
        for row in plots:
            if int(row[0])%3!=0:
                continue
            PX.append(int(row[0]))
            PY.append(float(row[1]))   
    with open('./upair/domain/distance.txt','r') as file:
        plots = csv.reader(file,delimiter = ',')
        for row in plots:
            if int(row[0])%3!=0:
                continue
            UX.append(int(row[0]))
            UY.append(float(row[1]))
    for x,y in zip(PX,PY):
        if x==33:
            plt.scatter(x,y,c='b',s=4)
            plt.text(x-7,y-0.0,(x,y),ha='center',va='bottom',fontsize=12)
    for x,y in zip(UX,UY):
        if x==33:
            plt.scatter(x,y,c='r',s=4)
            plt.text(x+7,y-0.05,(x,y),ha='center',va='bottom',fontsize=12)
    plt.axline([33,0.8],[33,0.311],color='y',linestyle='--')
    plt.plot(PX,PY,'bs-',label = '别名IP')
    plt.plot(UX,UY,'ro--',label = '非别名IP')
    plt.xlabel('域名编辑距离',fontsize=16)
    plt.ylabel('累积分布函数',fontsize=16)
    plt.tick_params(labelsize=13)
    plt.legend(loc='lower right',fontsize=13)
    plt.savefig("./pac/pic6."+format,bbox_inches='tight',format=format,dpi=dpi)

def netsec():
    plt.clf()
    plt.rcParams['font.sans-serif']=['SimHei']
    PX = []
    PY = []
    UX = []
    UY = []
    with open('./pair/netsec/netsec.csv','r') as file:
        plots = csv.reader(file,delimiter=',')
        for row in plots:
            if int(row[0])%3!=0 or int(row[0])==0:
                continue
            PX.append(int(row[0]))
            PY.append(float(row[1]))   
    with open('./upair/netsec/netsec.csv','r') as file:
        plots = csv.reader(file,delimiter = ',')
        for row in plots:
            if int(row[0])%3!=0 or int(row[0])==0:
                continue
            UX.append(int(row[0]))
            UY.append(float(row[1]))
    for x,y in zip(PX,PY):
        if x==24:
            plt.scatter(x,y,c='b',s=4)
            plt.text(x+3,y-0.03,(x,y),ha='center',va='bottom',fontsize=10)
        if x==3:
            plt.scatter(x,y,c='r',s=4)
            plt.text(x+3,y,(x,y),ha='center',va='bottom',fontsize=10)
    for x,y in zip(UX,UY):
        if x==24:
            plt.scatter(x,y,c='r',s=4)
            plt.text(x+3.5,y-0.005,(x,y),ha='center',va='bottom',fontsize=12)
        if x == 3:
            plt.scatter(x, y, c='r', s=4)
            plt.text(x + 2, y+0.03, (x, y), ha='center', va='bottom', fontsize=12)
    plt.axline([24,0.8],[24,0.311],color='y',linestyle='--')
    plt.plot(PX,PY,'bs-',label = '别名IP')
    plt.plot(UX,UY,'ro--',label = '非别名IP')
    plt.xlabel('IP空间距离',fontsize=16)
    plt.ylabel('累积分布函数',fontsize=16)
    plt.legend(loc='best',fontsize=13)
    plt.tick_params(labelsize=13)
    plt.savefig("./pac/pic7."+format,bbox_inches='tight',format=format,dpi=dpi)

def diff():
    plt.clf()
    plt.rcParams['font.sans-serif']=['SimHei']
    fig=plt.figure(figsize=(12,14))
    fig.subplots_adjust(hspace=0.3,wspace=0.3)
    LPX = []
    LPY = []
    LUX = []
    LUY = []
    plt.subplot(2,1,1)
    with open('./pair/diff/length.csv','r') as file:
        plots = csv.reader(file,delimiter=',')
        for row in plots:
            if int(row[0])<13:
                LPX.append(int(row[0]))
                LPY.append(float(row[1]))   
        for x,y in zip(LPX,LPY):
            plt.scatter(x,y,c='b',s=4)
            if x==4 :
                plt.text(x-0.8,y-0.,(x,y),ha='center',va='bottom',fontsize=16)
    with open('./upair/diff/length.csv','r') as file:
        plots = csv.reader(file,delimiter = ',')
        for row in plots:
            if int(row[0])<13:
                LUX.append(int(row[0]))
                LUY.append(float(row[1]))
        for x,y in zip(LUX,LUY):
            plt.scatter(x,y,c='r',s=5)
            if x==4:
                plt.text(x+0.8,y-0.1,(x,y),ha='center',va='bottom',fontsize=16)
    plt.axline([4,0.8],[4,0],color='y',linestyle='--')
    plt.plot(LPX,LPY,'bs-',label = '别名IP')
    plt.plot(LUX,LUY,'ro--',label = '非别名IP')
    plt.xlabel(r'路径长度差值',fontsize=22)
    plt.ylabel('累积分布函数',fontsize=22)
    plt.tick_params(labelsize=20)
    plt.title('a.路径长度差值',fontsize=20)
    plt.legend(loc='lower right',fontsize=20)
    # plt.yticks(fontsize=13)
    plt.subplot(2,1,2)
    DPX = []
    DPY = []
    DUX = []
    DUY = []
    with open('./pair/diff/direct.csv','r') as file:
        plots = csv.reader(file,delimiter=',')
        for row in plots:
            if int(row[0])<15:
                DPX.append(int(row[0]))
                DPY.append(float(row[1]))   
    with open('./upair/diff/direct.csv','r') as file: 
        plots = csv.reader(file,delimiter = ',')
        for row in plots:
            if int(row[0])<15:
                DUX.append(int(row[0]))
                DUY.append(float(row[1]))
    for x,y in zip(DPX[1:],DPY[1:]):
        plt.scatter(x,y,c='b',s=5)
        if  x==8:
            plt.text(x-0.8,y-0,(x,y),ha='center',va='bottom',fontsize=16)
        # if x==3 or x==4:
        #     plt.text(x+0.5,y-0.05,(x,y),ha='center',va='bottom',fontsize=8)
    # plt.axis([0.5, 14, 0, 1])
    for x,y in zip(DUX[1:],DUY[1:]):
        plt.scatter(x,y,c='r',s=5)
        if  x==8:
            plt.text(x+0.8,y-0.1,(x,y),ha='center',va='bottom',fontsize=16)
    plt.plot(DPX[1:],DPY[1:],'bs-',label = '别名IP')
    plt.plot(DUX[1:],DUY[1:],'ro--',label = '非别名IP')
    plt.axline([8,0.8],[8,0],color='y',linestyle='--')
    plt.xlabel(r'路径方向差值',fontsize=22)
    plt.ylabel('累积分布函数',fontsize=22)
    plt.legend(loc='lower right',fontsize=20)
    plt.tick_params(labelsize=20)
    plt.title('b.路径方向差值',fontsize=20)
    plt.savefig("./pac/pic3."+format,bbox_inches='tight',format=format,dpi=dpi)
    # plt.show()

# def compare():
#     # x_label=['0',r'$\regular{10^0}$',r'$\regular{10^1}$',r'$\regular{10^2}$',r'$\regular{10^3}$',r'$\regular{10^4}$',r'$\regular{10^5}$',r'$\regular{10^6}$',r'$\regular{10^7}$']
#     X = [10000,100000,500000,1000000,2000000]
#     X_num = [0,1,10,100,1000,10000,100000,1000000,10000000]
#     X_fake = [4,5,5.69897,6,6.30102]
#     x_fake_label=['0',r'$\regular{10^0}$',r'$\regular{10^1}$',r'$\regular{10^2}$',r'$\regular{10^3}$',r'$\regular{10^4}$',r'$\regular{10^5}$',r'$\regular{10^6}$',r'$\regular{10^7}$']
#     X_index = [0,1,2,3,4,5,6,7,8]
#     y_index = [0,1,2,3,4,5]
#     y_label = ['0',r'$\regular{10^0}$',r'$\regular{10^1}$',r'$\regular{10^2}$',r'$\regular{10^3}$',r'$\regular{10^4}$']
#     # X_index = ['1e4','1e5','5e5','1e6','2e6']
#     TY = [67.8,457.2,715.8,5227.8,10947.6]
#
#     MY = [750,90000,2202264,2202264,2202264]
#     AY = [1.5,14.9,74.4,148.8,297]
#     CY = [2.22,30.71,348.6,1162.2,4115.4]
#     for i in range(0,5):
#         TY[i]=math.log(TY[i]/60,10)
#         MY[i]=math.log(MY[i]/60,10)
#         AY[i]=math.log(AY[i]/60,10)
#         CY[i]=math.log(CY[i]/60,10)
#     # _ = plt.xticks(X, X_index)
#     plt.axis([0, 6, 0, 5])
#     # plt.xscale('symlog')
#     # plt.yscale('symlog')
#     # ax=plt.gca()
#     # ax.margins(x=10000)
#     plt.plot(X_fake,TY,label = 'TreeNet',color='red',linestyle='-',marker='s')
#     plt.plot(X_fake,MY,label = 'MLAR',color='green',linestyle='--',marker='^')
#     plt.plot(X_fake,AY,label = 'APPLE',color='black',linestyle='-.',marker='x')
#     plt.plot(X_fake,CY,label = 'AliasClassifier',color='blue',linestyle=':',marker='o')
#     # ax.spines['left'].set_position(('data', 0))
#     plt.xlabel('NO.ip')
#     plt.ylabel('T/minute')
#     plt.xticks(X_index, x_fake_label)
#     plt.yticks(y_index,y_label)
#     plt.legend(loc='best')
#     plt.savefig("./pac/compare",bbox_inches='tight',format=format)


def fbar():
    plt.clf()
    plt.rcParams['font.sans-serif']=['SimHei']
    plt.figure(figsize=(7,5))
    #pre,rec,f1,f0.5,f2
    ys = [0.9989,0.0150,0.0267,0.0706,0.0186,0.9284,0.0885,0.1616,0.3203,0.1080,1,0.0022,0.0044,0.011,0.0028,0.9486,0.4045,0.5670,0.7475,0.4569]
    treent_y = ys[:2]
    mlar = ys[5:7]
    apple = ys[10:12]
    aliasC = ys[15:17]
    x_label = ['TreeNet','MLAR','APPLE','AliasClassifier']
    width = 0.8
    xsec = [2,4,6,8]
    x2 = [xsec[0] + width / 2, xsec[1] + width / 2, xsec[2] + width / 2,xsec[3]+width/2]
    x1 = [xsec[0] - width / 2, xsec[1] - width / 2, xsec[2] - width / 2,xsec[3]-width/2]
    y1 = [treent_y[0],mlar[0],apple[0],aliasC[0]]
    y2 = [treent_y[1],mlar[1],apple[1],aliasC[1]]
    # xticks = 3
    plt.bar(x1, y1, width=width, hatch="", label='精度', color='white',edgecolor='black')
    plt.bar(x2, y2, width=width, hatch="/", label='召回率', color='grey',edgecolor='black')
    plt.xlabel('解析器', fontsize=16, labelpad=8)
    plt.ylabel('得分', fontsize=16, labelpad=8)
    for a, b in zip(x1, y1):
        # if a==2-width/2:
        #   plt.text(a, b, '%.4f' % b, ha='center', va='bottom', fontsize=10)
        # else:
          plt.text(a, b, '%.4f' % b, ha='center', va='bottom', fontsize=9)
    for a, b in zip(x2, y2):
        plt.text(a, b, '%.4f' % b, ha='center', va='bottom', fontsize=9)
    plt.xticks(xsec, x_label)
    # plt.axis([0, 9, 0, 1.19])
    plt.legend(loc='best')
    plt.margins(0.11)
    plt.savefig("./pac/pic9."+format, bbox_inches='tight',format=format,dpi=dpi)
    # plt.show()

def fparser5bar():
    plt.clf()
    plt.rcParams['font.sans-serif']=['SimHei']
    fig=plt.figure(figsize=(16,8))
    #pre,rec,f1,f0.5,f2
    ys = [0.9989,0.0150,0.0267,0.0706,0.0186,0.9284,0.0885,0.1616,0.3203,0.1080,1,0.0022,0.0044,0.011,0.0028,0.9486,0.4045,0.5670,0.7475,0.4569]
    treent_y = ys[2:5]
    mlar = ys[7:10]
    apple = ys[12:15]
    aliasC = ys[17:20]
    x_label = ['TreeNet','MLAR','APPLE','AliasClassifier']
    width = 1.7
    # xsec = [2,4,6,8]
    base=6
    xsec=[1*base,2*base,3*base,4*base]
    x1 = []
    x2 = []
    x3 = []
    # x4 = []
    # x5 = []
    for x in xsec:
        x1.append(x-width)
        x2.append(x)
        x3.append(x+width)
        # x4.append(x+width)
        # x5.append(x+2*width)
    # y1 = []
    # y2 = []
    # y2 = [treent_y[3],mlar[3],apple[3],aliasC[3]]
    # y3 = [treent_y[4],mlar[4],apple[4],aliasC[4]]
    y1 = [treent_y[1],mlar[1],apple[1],aliasC[1]]
    y2 = [treent_y[0],mlar[0],apple[0],aliasC[0]]
    y3 = [treent_y[2],mlar[2],apple[2],aliasC[2]]
    # for i in range(2,5):

    # xticks = 3
    plt.bar(x1, y1, width=width, hatch="x", label=r'$\regular{F_{0.5}}$', color='white',edgecolor='black')
    plt.bar(x2, y2, width=width, hatch="", label=r'$\regular{F_{1}}$', color='grey',edgecolor='black')
    plt.bar(x3,y3,width=width,hatch='/',label=r'$\regular{F_{2}}$',color='cyan',edgecolor='black')
    # plt.bar(x4,y4,width=width,hatch='/',label='f0.5',color = 'wheat',edgecolor='black')
    # plt.bar(x5,y5,width=width,hatch='.',label='f2',color='magenta',edgecolor='black')
    plt.xlabel('解析器', fontsize=22, labelpad=8)
    plt.ylabel('得分', fontsize=22, labelpad=8)
    for a, b in zip(x1, y1):
          plt.text(a, b, '%.3f' % b, ha='center', va='bottom', fontsize=17)
    for a, b in zip(x2, y2):
        plt.text(a, b, '%.3f' % b, ha='center', va='bottom', fontsize=17)

    for a, b in zip(x3, y3):
        plt.text(a, b, '%.3f' % b, ha='center', va='bottom', fontsize=17)

    plt.xticks(xsec, x_label)
    plt.tick_params(labelsize=20)
    # plt.axis([0, 9, 0, 1.19])
    plt.legend(loc='upper left',fontsize=20)
    # plt.margins(0.2)
    # plt.margins(x=0.0,y=0.2)
    plt.axis([0+2,base*5-2, 0, 0.8])
    plt.savefig("./pac/pic10."+format, bbox_inches='tight',format=format,dpi=dpi)
def compare_easy():
    plt.clf()
    plt.rcParams['font.sans-serif']=['SimHei']
    plt.figure(figsize=(8,5))
    # x_label=['0',r'$\regular{10^0}$',r'$\regular{10^1}$',r'$\regular{10^2}$',r'$\regular{10^3}$',r'$\regular{10^4}$',r'$\regular{10^5}$',r'$\regular{10^6}$',r'$\regular{10^7}$']
    X = [10000,100000,500000,1000000,2000000]
    y_label = ['0',r'$\regular{10^0}$',r'$\regular{10^1}$',r'$\regular{10^2}$',r'$\regular{10^3}$']
    TY = [67.8,457.2,715.8,5227.8,10947.6]
    MY = [750,90000,2202264,2202264,2202264]
    AY = [1.5,14.9,74.4,148.8,297]
    CY = [2.22,30.71,348.6,1162.2,4115.4]
    y_index = [0,1,2,3,4]
    for i in range(0,5):
        TY[i]=TY[i]/60
        MY[i]=MY[i]/60
        AY[i]=AY[i]/60
        CY[i]=CY[i]/60
    print(MY)
    print(math.log(MY[1],10))
    for i in range(0,5):
        if TY[i]>1:
            TY[i]=math.log(TY[i],10)+1
        if MY[i]>1:
            MY[i]=math.log(MY[i],10)+1
        if AY[i]>1:
            AY[i]=math.log(AY[i],10)+1
        if CY[i]>1:
            CY[i]=math.log(CY[i],10)+1
    # _ = plt.xticks(X, X_index)
    plt.axis([8000, 3000000, -0.3, 5])
    # plt.axis([])
    plt.xscale('symlog')
    # plt.yscale('symlog')
    # ax=plt.gca()
    # ax.margins(x=10000)
    plt.plot(X,TY,label = 'TreeNet',color='red',linestyle='-',marker='s')
    plt.plot(X,MY,label = 'MLAR',color='green',linestyle='--',marker='^')
    plt.plot(X,AY,label = 'APPLE',color='black',linestyle='-.',marker='x')
    plt.plot(X,CY,label = 'AliasClassifier',color='blue',linestyle=':',marker='o')
    # ax.spines['left'].set_position(('data', 0))
    plt.xlabel('IP数量',fontsize=16,labelpad=10)
    plt.ylabel('时间/小时',fontsize=16,labelpad=10)
    # plt.xticks(X_index, x_fake_label)
    plt.yticks(y_index,y_label)
    plt.legend(loc='best',fontsize=12)
    plt.savefig("./pac/pic11."+format,bbox_inches='tight',format=format,dpi=dpi)


# editDistance()
# factor()
# rtt()
# ttl()
# domain()
# netsec()
# diff()
# compare()
# fbar()
fparser5bar()
compare_easy()